package main

import (
	"context"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"rusprofile-fetcher/internal/rpc_server"
	"strings"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

//go:embed static
var content embed.FS

func main() {
	// адрес для HTTP запросов
	proxyAddr := ":8000"

	// адрес для gRPC запросов
	serviceAddr := "127.0.0.1:8001"

	go gRPCService(serviceAddr)
	HTTPProxy(proxyAddr, serviceAddr)
}

func gRPCService(addr string) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	server := grpc.NewServer()

	rpc_server.RegisterOrgInfoServiceServer(server, newServer())

	fmt.Println("начинаем слушать gRPC сервер")
	server.Serve(l)
}

func HTTPProxy(proxyAddr, serviceAddr string) {
	grpcConn, err := grpc.Dial(serviceAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer grpcConn.Close()

	grpcMux := runtime.NewServeMux()

	if err := rpc_server.RegisterOrgInfoServiceHandler(context.Background(), grpcMux, grpcConn); err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()

	mux.Handle("/inn/", grpcMux)
	mux.HandleFunc("/openapi.json", swaggerHandler)

	fsys, _ := fs.Sub(content, "static/swagger-ui")
	mux.Handle("/", http.StripPrefix("/", http.FileServer(http.FS(fsys))))

	fmt.Println("начинаем слушать HTTP сервер")
	log.Fatal(http.ListenAndServe(proxyAddr, mux))
}

var apiURL = "https://www.rusprofile.ru/ajax.php?query=%s&action=search&cacheKey=%.12f"

type server struct {
	rpc_server.UnimplementedOrgInfoServiceServer
}

func newServer() *server {
	return &server{}
}

func (s *server) Fetch(ctx context.Context, _in *rpc_server.Request) (*rpc_server.Response, error) {
	if len(_in.GetINN()) != 10 && len(_in.GetINN()) != 12 && len(_in.GetINN()) != 5 {
		return nil, errors.New("ИНН должен быть 12 цифр для ФЛ и ИП, 10 цифр для ЮЛ или 5 для иностранных ЮЛ")
	}

	rand.Seed(time.Now().Unix())

	resp, err := http.Get(fmt.Sprintf(apiURL, _in.GetINN(), rand.Float64()))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	r := struct {
		UL []struct {
			INN  string `json:"inn"`
			Name string `json:"name"`
			CEO  string `json:"ceo_name"`
			OGRN string `json:"ogrn"` // КПП нет с запроса, а парсить HTML, мне кажется не лучшее решение
		} `json:"ul"`
	}{}

	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}

	io.Copy(io.Discard, resp.Body)

	if len(r.UL) == 0 {
		return nil, errors.New("организация не найдена")
	}

	r.UL[0].INN = strings.Trim(r.UL[0].INN, "~!")

	return &rpc_server.Response{
		INN:      r.UL[0].INN,
		OrgName:  r.UL[0].Name,
		Director: r.UL[0].CEO,
		OGRN:     r.UL[0].OGRN,
	}, nil
}

func swaggerHandler(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadFile("internal/rpc_server/server.swagger.json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":  err.Error(),
			"result": nil,
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(b)
}
