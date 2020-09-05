package main

import (
	"math/rand"
	"time"

	pb "github.com/mmontes11/go-grpc-routes/pb"
)

var validPoint = &pb.Point{Latitude: 409146138, Longitude: -746188906}

var invalidPoint = &pb.Point{Latitude: 0, Longitude: 0}

var rect = &pb.Rectangle{
	Lo: &pb.Point{Latitude: 410000000, Longitude: -740000000},
	Hi: &pb.Point{Latitude: 415000000, Longitude: -745000000},
}

func randomPoints() []*pb.Point {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	pointCount := int(r.Int31n(100)) + 2
	var points []*pb.Point
	for i := 0; i < pointCount; i++ {
		points = append(points, randomPoint(r))
	}
	return points
}

func randomPoint(r *rand.Rand) *pb.Point {
	lat := (r.Int31n(180) - 90) * 1e7
	long := (r.Int31n(360) - 180) * 1e7
	return &pb.Point{Latitude: lat, Longitude: long}
}

func notes() []*pb.RouteNote {
	return []*pb.RouteNote{
		{Location: &pb.Point{Latitude: 0, Longitude: 1}, Message: "First message"},
		{Location: &pb.Point{Latitude: 0, Longitude: 2}, Message: "Second message"},
		{Location: &pb.Point{Latitude: 0, Longitude: 3}, Message: "Third message"},
		{Location: &pb.Point{Latitude: 0, Longitude: 1}, Message: "Fourth message"},
		{Location: &pb.Point{Latitude: 0, Longitude: 2}, Message: "Fifth message"},
		{Location: &pb.Point{Latitude: 0, Longitude: 3}, Message: "Sixth message"},
	}
}
