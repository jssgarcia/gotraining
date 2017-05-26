package lorem

import golorem "github.com/drhodes/golorem"

type Service interface {
	Word(min,max int) string
	Sentence(min,max int) string
	Paragraph(min,max int) string
}

type LoremService struct {

}

func (LoremService) Word(min,max int) string {
	return golorem.Word(min,max)
}

func (LoremService) Sentence(min,max int) string {
	return golorem.Sentence(min,max)
}

func (LoremService) Paragraph(min,max int) string {
	return golorem.Paragraph(min,max)
}