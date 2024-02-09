package workflow

import (
	"github.com/gin-gonic/gin"
	"github.com/remiges-tech/alya/service"
	"github.com/remiges-tech/alya/wscutils"
)

type Point struct {
	X int
	Y int
}

func Test(c *gin.Context, s *service.Service) {

	// nestedArray := [][]int{
	// 	{1, 2, 3},
	// 	{4, 5, 6},
	// 	{7, 8, 9},
	// }
	arrayOfArrays := [][]Point{
		{Point{1, 1}, Point{1, 2}, Point{1, 3}, Point{1, 4}},
		{Point{2, 1}, Point{2, 2}, Point{2, 3}, Point{2, 4}},
		{Point{3, 1}, Point{3, 2}, Point{3, 3}, Point{3, 4}},
	}

	wscutils.SendSuccessResponse(c, &wscutils.Response{
		Data: arrayOfArrays,
	})
}
