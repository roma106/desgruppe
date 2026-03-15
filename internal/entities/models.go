package entities

import "time"

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	SectionID   int     `json:"sectionId"`
	ProducerID  int     `json:"producerId"`
	DesignerID  int     `json:"designerId"`
	Size        string  `json:"size"`
	Available   bool    `json:"available"`
	Price       float64 `json:"price"`
	OnSale      bool    `json:"onSale"`
	Sale        float64 `json:"sale"`
	Photo       []int   `json:"photo"`
	Description string  `json:"description"`
	Colors      []int   `json:"colors"`
	Position    int     `json:"position"`
	Slug        string  `json:"slug"`
	FreeForm    any     `json:"freeForm"`
	Best        bool    `json:"best"`
}

type ProductPhoto struct {
	ID        int    `json:"id"`
	ProductID string `json:"productId"`
	Photo     []int  `json:"photo"`
	Position  int    `json:"position"`
}

type Color struct {
	ID        int    `json:"id" db:"ID"`
	Name      string `json:"name" db:"Name"`
	Code      string `json:"code" db:"Code"`
	Picture   []int  `json:"picture" db:"Picture"`
	Position  int    `json:"position" db:"Position"`
	IsCode    bool   `json:"codeRadio"`
	IsPicture bool   `json:"pictureRadio"`
}

type TemplColor struct {
	ID        int    `json:"id" db:"ID"`
	Name      string `json:"name" db:"Name"`
	Code      string `json:"code" db:"Code"`
	Picture   []int  `json:"picture" db:"Picture"`
	Position  int    `json:"position" db:"Position"`
	IsCode    bool   `json:"codeRadio"`
	IsPicture bool   `json:"pictureRadio"`
}

type Designer struct {
	ID          int    `json:"id" db:"ID"`
	Name        string `json:"name" db:"Name"`
	Description string `json:"description" db:"Description"`
	Picture     []int  `json:"picture" db:"Picture"`
	Position    int    `json:"position" db:"Position"`
	Slug        string `json:"slug"`
}

type Producer struct {
	ID          int    `json:"id" db:"ID"`
	Name        string `json:"name" db:"Name"`
	Description string `json:"description" db:"Description"`
	Picture     []int  `json:"picture" db:"Picture"`
	Position    int    `json:"position" db:"Position"`
	Slug        string `json:"slug"`
}

type Section struct {
	ID   int    `json:"id" db:"ID"`
	Name string `json:"name" db:"Name"`
	Type string `json:"type" db:"Type"`
}

type Settings struct {
	ID           int     `json:"id" db:"ID"`
	ExchangeRate float64 `json:"exchangeRate" db:"ExchangeRate"`
	Email        string  `json:"email" db:"Email"`
}

type CartProduct struct {
	ID           int `json:"id" db:"ID"`
	ProductID    int `json:"productID" db:"ProductIDs"`
	ColorID      int `json:"colorID" db:"ColorIDs"`
	QuantitiesID int `json:"quantitiesID" db:"QuantitiesIDs"`
	Index        int `json:"index"`
}
type Cart struct {
	ID            int   `json:"id" db:"ID"`
	ProductIDs    []int `json:"productIDs" db:"ProductIDs"`
	ColorIDs      []int `json:"colorIDs" db:"ColorIDs"`
	QuantitiesIDs []int `json:"quantitiesIDs" db:"QuantitiesIDs"`
}

type Order struct {
	ID      int       `json:"id" db:"ID"`
	Date    time.Time `json:"date" db:"Date"`
	Name    string    `json:"name" db:"Name"`
	Phone   string    `json:"phone" db:"Phone"`
	Email   string    `json:"email" db:"Email"`
	Comment string    `json:"comment" db:"Comment"`
	CartID  string    `json:"cartId" db:"CartID"`
	Seen    bool      `db:"Seen"`
}

type Question struct {
	ID      int       `json:"id" db:"ID"`
	Date    time.Time `json:"date" db:"Date"`
	Name    string    `json:"name" db:"Name"`
	Phone   string    `json:"phone" db:"Phone"`
	Email   string    `json:"email" db:"Email"`
	Message string    `json:"message" db:"Message"`
	Seen    bool      `json:"seen" db:"Seen"`
}

type Sale struct {
	ID        int    `json:"id" db:"ID"`
	ProductID string `json:"productId" db:"ProductID"`
	Position  int    `json:"position" db:"Position"`
}
