package blog

// Posts is the only publication registry. Register real posts here; keep the
// example disabled as a reference for future authors.
func Posts() []Post {
	return []Post{
		// ExamplePost(), // Reference only; do not publish.
	}
}

var exampleHeadings = []Heading{
	{ID: "prepoznajte-proces", Title: "1. Prepoznajte dobar prvi proces"},
	{ID: "usporedite-pristupe", Title: "2. Usporedite pristupe"},
	{ID: "izmjerite-polaziste", Title: "3. Izmjerite polazište"},
	{ID: "zadrzite-kontrolu", Title: "4. Zadržite čovjeka u procesu"},
	{ID: "dokumentirajte-pravila", Title: "5. Dokumentirajte pravila"},
	{ID: "prvi-korak", Title: "Prvi korak: mali, provjerljiv test"},
	{ID: "cesta-pitanja", Title: "Česta pitanja"},
	{ID: "izvori", Title: "Izvori i bilješke"},
}

func ExamplePost() Post {
	return Post{
		Slug:        "primjer-blog-objave",
		Title:       "Kako prepoznati proces koji vrijedi automatizirati",
		Description: "Praktičan okvir za odabir prvog poslovnog procesa za automatizaciju: ponavljanje, jasna pravila, mjerenje vremena i ljudska kontrola.",
		Excerpt:     "Odaberite mali, ponovljiv proces. Izmjerite ručni rad, definirajte iznimke i provjerite rezultat prije širenja rješenja.",
		PublishedAt: "2026-10-24T00:00:00+02:00",
		Status:      Example,
		Topics:      []string{"Automatizacija", "Poslovni procesi"},
		Hero: Image{
			Src: "/assets/blog/example/hero.jpg", Alt: "Ilustracija obrade upita: provjera podataka, nacrt ponude, ljudska potvrda i obrada iznimke",
			Width: 2000, Height: 800,
			Variants: []Source{{Src: "/assets/blog/example/hero-800.jpg", Width: 800, Height: 320}, {Src: "/assets/blog/example/hero-1200.jpg", Width: 1200, Height: 480}},
		},
		Social: Image{
			Src: "/assets/blog/example/social.jpg", Alt: "Od upita do provjerene ponude — primjer automatizacije poslovnog procesa",
			Width: 1200, Height: 630,
		},
		Headings: exampleHeadings,
		Body:     exampleBody(),
	}
}
