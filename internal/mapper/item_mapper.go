package mapper

import (
	"fmt"
	"html/template"
	"image"
	_ "image/png"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/iambaangkok/Card-Maker/internal/entity"
	"github.com/iambaangkok/Card-Maker/internal/reader"
)

type ItemMapper interface {
	Map(csvFile reader.CSVFile) []entity.Item
}

type ItemMapperImpl struct {
	ExistingEffects map[string]entity.Effect
}

func (w ItemMapperImpl) Map(csvFile reader.CSVFile) []entity.Item {
	
	var items []entity.Item
	
	for row, line := range csvFile.Records {
		log.Println("mapping line", row, line)

		usageLimit, err := strconv.Atoi(line[1])
		if err != nil { log.Fatal("usage limit must be int") }
		price, err := strconv.Atoi(line[2])
		if err != nil { log.Fatal("price must be int") }

		effectStrs := strings.Split(line[3], "/")
		var effects []entity.Effect
		for _, effectStr := range effectStrs {
			val := strings.Split(effectStr, ":")
			effectName := val[0]
			if effectName == "-" {
				continue
			}

			effect, exists := w.ExistingEffects[effectName]
			if !exists { log.Fatal("invalid effect name") }

			if effect.HasLevel {
				if len(val) == 1 { 
					effect.Level = 1
				}else{
					level, err := strconv.Atoi(val[1])
					if err != nil {
						log.Fatal("level must be int")
					}
					effect.Level = level
				}
			}
			effects = append(effects, effect)
		}
		
		// ImgTag
		imgTagStr := `<img src=%s class="FrameC1" style=%s />`
		// imageServerPath
		imageServerPath := fmt.Sprintf(`"http://localhost:8081/static/img/items/%s.png"`, line[0])
		// imageStyle
		filePath := filepath.Join("./internal/template/img/items/", line[0] + ".png")
		log.Println(filePath)
		reader, err := os.Open(filePath)
		if err != nil {
			log.Fatal(err)
		}
		defer reader.Close()
		im, _, err := image.DecodeConfig(reader)
		if err != nil {
			log.Fatal(err)
		}
		divisor := 6
		imageStyle := fmt.Sprintf(`"width: auto; height: 100%%; max-width: %dpx; max-height: %dpx;
		margin:auto; align-self: center; flex: 0 1 0; object-fit: contain;"`,
		im.Width/divisor, im.Height/divisor)
		imgTagStr = fmt.Sprintf(imgTagStr, imageServerPath, imageStyle)
		imgTag := template.HTML(imgTagStr)

		items = append(items, 
		entity.Item{
			Name: line[0],
			UsageLimit: usageLimit,
			Price: price,
			Effects: effects,
			ImgTag: imgTag,
		})
	}
	return items
}