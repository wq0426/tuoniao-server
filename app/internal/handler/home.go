package handler

import (
	"log"

	"github.com/gin-gonic/gin"

	"app/internal/common"
	pb "app/internal/grpc"
	"app/internal/service"
)

type HomeHandler struct {
	*Handler
	pb.UnimplementedPushMessageServiceServer
	EventChannel        map[string]chan *pb.PushMessageResponse // 用于接收事件消息
	EventHistoryChannel map[string]pb.PushMessageResponse       // 未处理的历史事件消息
	Streams             map[string]pb.PushMessageService_StreamMessagesServer
}

func NewHomeHandler(
	handler *Handler,
) *HomeHandler {
	return &HomeHandler{
		Handler: handler,
	}
}

// 事件触发器，用于模拟外部事件
func (s *HomeHandler) TriggerEvent(ctx *gin.Context, event *pb.PushMessageResponse) {
	userId := service.GetUserIdFromCtx(ctx)
	if _, ok := s.EventChannel[userId]; !ok {
		log.Printf("Error s.EventChannel not exist userId : %s", userId)
		s.EventHistoryChannel[userId] = *event
	} else {
		// 向事件通道发送事件
		s.EventChannel[userId] <- event
	}
}

// StreamMessages 方法，服务端通过事件推送消息
func (s *HomeHandler) StreamMessages(stream pb.PushMessageService_StreamMessagesServer) error {
	// 服务端接收客户端的请求并响应
	for {
		// 接收客户端发送的消息
		req, err := stream.Recv()
		if err != nil || len(req.UserId) == 0 {
			log.Printf("Error receiving message: %v", err)
			return err
		}
		// 是否属于系统推送
		if req.UserId == "system" {
			if _, ok := s.EventChannel[req.UserId]; !ok {
				s.Streams[req.UserId] = stream
				s.EventChannel[req.UserId] = make(chan *pb.PushMessageResponse, 10)
				go func() {
					// 主动事件触发：模拟服务端定时推送消息
					quit := false
					for {
						if quit {
							break
						}
						select {
						case <-s.Streams[req.UserId].Context().Done():
							s.Clear(req.UserId)
							quit = true
							break
						case res, ok := <-s.EventChannel[req.UserId]:
							if !ok {
								s.Clear(req.UserId)
								quit = true
								break
							}
							err = stream.Send(res)
							if err != nil {
								log.Printf("Error sending event to client: %v", err)
							}
						}
					}
				}()
			}
			if len(req.To) > 0 {
				if _, ok := s.EventChannel[req.To]; ok {
					s.EventChannel[req.To] <- service.Transfer(common.MODULE_TYPE_ESCORT, req.Data)
				}
			}
		}
		// 是否存在用户状态变更的历史事件消息
		if _, ok := s.EventHistoryChannel[req.UserId]; ok {
			historyChannel := s.EventHistoryChannel[req.UserId]
			err = stream.Send(&historyChannel)
			if err != nil {
				log.Printf("Error sending EventHistoryChannel event to client: %v, userId: %s", err, req.UserId)
				return err
			}
			delete(s.EventHistoryChannel, req.UserId)
		}
		if _, ok := s.EventChannel[req.UserId]; !ok {
			s.Streams[req.UserId] = stream
			s.EventChannel[req.UserId] = make(chan *pb.PushMessageResponse, 5)
			go func() {
				// 主动事件触发：模拟服务端定时推送消息
				quit := false
				for {
					if quit {
						break
					}
					select {
					case <-s.Streams[req.UserId].Context().Done():
						s.Clear(req.UserId)
						quit = true
						break
					case res, ok := <-s.EventChannel[req.UserId]:
						if !ok {
							s.Clear(req.UserId)
							quit = true
							break
						}
						err = stream.Send(res)
						if err != nil {
							log.Printf("Error sending event to client: %v", err)
						}
					}
				}
			}()
		}
		// rsponse, err := s.homeService.GetHome(stream.Context(), req.UserId)
		// if err != nil {
		// 	log.Printf("Error homeService.GetHome message: %v", err)
		// 	return err
		// }
		// err = stream.Send(service.Transfer(common.MODULE_TYPE_USER_INFO, rsponse))
		// if err != nil {
		// 	log.Printf("Error sending event to client: %v", err)
		// 	return err
		// }
	}
}

func (s *HomeHandler) Clear(userId string) {
	delete(s.Streams, userId)
	delete(s.EventChannel, userId)
}
