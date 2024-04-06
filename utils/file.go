package utils

import (
	"errors"
	"fmt"
	"github.com/gosimple/slug"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"
)

const DefaultPathAsset = "./public/assets/"

func HandleMultipleFile(ctx *fiber.Ctx) error {
	form, err := ctx.MultipartForm()
	if err != nil {
		log.Println("Error read multipart form request, Error :", err.Error())
	}

	files := form.File["product_images"]
	productSlug := slug.Make(form.Value["name"][0])
	var filenames []string
	for idx, file := range files {
		if file != nil {
			err = CheckContentType(file, "image/jpg", "image/png", "image/jpeg")
			if err != nil {
				return ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
					"message": err.Error(),
				})
			}

			extFile := filepath.Ext(file.Filename)
			filename := fmt.Sprintf("%d-%s-%d%s", idx, productSlug, time.Now().UnixMilli(), extFile)

			err = ctx.SaveFile(file, DefaultPathAsset+filename)
			if err != nil {
				log.Println("Failed to store image")
			}

			filenames = append(filenames, filename)
		} else {
			log.Println("Nothing file to upload")
		}
	}

	ctx.Locals("filenames", filenames)

	return ctx.Next()
}

func HandleRemoveFile(filename string, path ...string) error {
	if len(path) > 0 {
		err := os.Remove(path[0] + filename)
		if err != nil {
			log.Println("Failed to remove file")
			return err
		}
	} else {
		err := os.Remove(DefaultPathAsset + filename)
		if err != nil {
			log.Println("Failed to remove file")
			return err
		}
	}

	return nil
}

func CheckContentType(file *multipart.FileHeader, contentTypes ...string) error {
	if len(contentTypes) > 0 {
		for _, contentType := range contentTypes {
			typeFile := file.Header.Get("content-type")
			if typeFile == contentType {
				return nil
			}
		}

		return errors.New("only allowed png/jpg/jpeg file!")
	} else {
		return errors.New("FILE NOT FOUND TO CHECKING")
	}
}
