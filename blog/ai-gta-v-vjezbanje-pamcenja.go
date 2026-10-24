package blog

var gtaVMemoryHeadings = []Heading{
	{ID: "lose-se-snalazim-u-prostoru", Title: "Loše se snalazim u prostoru"},
	{ID: "ucenje-uz-ai", Title: "Kako se u ovu priču uklapa učenje uz AI?"},
	{ID: "mod-za-gta-v", Title: "Rezultat - mod za GTA V kao pomoć u učenju"},
	{ID: "koristi-ai-kreativno", Title: "AI je moćan alat, koristi ga kreativno"},
}

func GTAVMemoryPost() Post {
	return Post{
		Slug:        "ai-gta-v-vjezbanje-pamcenja",
		Title:       "Kako sam iskoristio AI i GTA V da bi vježbao pamćenje",
		Description: "Kako sam uz Claude Code izradio mod za GTA V koji mi pomaže vježbati pamćenje puteva i snalaženje u prostoru.",
		PublishedAt: "2026-10-24T16:43:50+02:00",
		Status:      Published,
		Topics:      []string{"AI", "Učenje", "Izrada softwarea"},
		// Original abstract route illustration, not a GTA V map or mod screenshot.
		Hero: Image{
			Src: "/assets/blog/ai-gta-v-vjezbanje-pamcenja/hero.jpg", Alt: "Apstraktna karta ulica s plavom rutom, skretanjima i orijentirima za vježbanje pamćenja puta",
			Width: 2000, Height: 800,
			Variants: []Source{
				{Src: "/assets/blog/ai-gta-v-vjezbanje-pamcenja/hero-800.jpg", Width: 800, Height: 320},
				{Src: "/assets/blog/ai-gta-v-vjezbanje-pamcenja/hero-1200.jpg", Width: 1200, Height: 480},
			},
		},
		Social: Image{
			Src: "/assets/blog/ai-gta-v-vjezbanje-pamcenja/social.jpg", Alt: "Plava ruta kroz apstraktnu kartu ulica — ilustracija pamćenja puteva, ne prikaz igre GTA V",
			Width: 1200, Height: 630,
		},
		Headings: gtaVMemoryHeadings,
		Body:     gtaVMemoryBody(),
	}
}
