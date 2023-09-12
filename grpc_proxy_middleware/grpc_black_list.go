package grpc_proxy_middleware

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/e421083458/go_gateway/dao"
	"github.com/e421083458/go_gateway/public"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

func GrpcBlackListMiddleware(serviceDetail *dao.ServiceDetail) func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		whileIpList := []string{}
		if serviceDetail.AccessControl.WhiteList != "" {
			whileIpList = strings.Split(serviceDetail.AccessControl.WhiteList, ",")
		}
		perrCtx, ok := peer.FromContext(ss.Context())
		if !ok {
			return errors.New("perr not found with context")
		}
		peerAddr := perrCtx.Addr.String()
		addrPos := strings.LastIndex(peerAddr, ":")
		clientIp := peerAddr[0:addrPos]
		blackIpList := []string{}
		if serviceDetail.AccessControl.BlackList != "" {
			blackIpList = strings.Split(serviceDetail.AccessControl.BlackList, ",")
		}
		if serviceDetail.AccessControl.OpenAuth == 1 && len(whileIpList) == 0 && len(blackIpList) > 0 {
			if public.InStringSlice(blackIpList, clientIp) {
				return errors.New(fmt.Sprintf("%s int black ip list", clientIp))
			}
		}
		if err := handler(srv, ss); err != nil {
			log.Printf("RPC failed with error %v\n", err)
			return err
		}
		return nil
	}
}
