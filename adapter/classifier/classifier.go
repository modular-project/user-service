package classifier

import (
	"context"
	"log"
	"mime/multipart"
	"net/http"
	"users-service/pkg"

	pb "github.com/modular-project/protobuffers/classification"
	"google.golang.org/grpc"
)

type ClassifierService struct {
	s pb.ClassImgServiceClient
}

func NewClassifierService(conn *grpc.ClientConn) ClassifierService {
	return ClassifierService{s: pb.NewClassImgServiceClient(conn)}
}

func (cls ClassifierService) ClassImg(ctx context.Context, h *multipart.FileHeader) (uint32, error) {
	f, err := h.Open()
	if err != nil {
		return 0, pkg.NewAppError("could not open file", err, http.StatusBadRequest)
	}
	log.Print(h.Size)
	if h.Size > 4000000 {
		return 0, pkg.NewAppError("file must be less than 4 MB", nil, http.StatusBadRequest)
	}
	b := make([]byte, h.Size)
	n, err := f.Read(b)
	if err != nil {
		return 0, pkg.NewAppError("failed to read file", err, http.StatusBadRequest)
	}
	log.Print("Readed: ", n)
	r, err := cls.s.ClassImg(ctx, &pb.RequestClassImg{Img: b})
	if err != nil {
		return 0, pkg.NewAppError("failed to classify image", err, http.StatusInternalServerError)
	}
	return r.Id, nil
}
