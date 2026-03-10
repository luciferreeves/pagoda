package districts

import "shrine/enums"

type District struct {
	Slug        enums.DistrictSlug
	Name        string
	Description string
	Background  string
	Foreground  string
	Detail      string
}

var All = []District{
	{enums.Arcadia, "Arcadia", "Video games, puzzles, toys", "#141c38", "#90c8ff", "#6898d0"},
	{enums.Arles, "Arles", "Drawings, photos, visual art", "#281420", "#ff98b0", "#d07088"},
	{enums.Hollywood, "Hollywood", "Cartoons, movies, western media", "#281808", "#ffb050", "#d08828"},
	{enums.Oxford, "Oxford", "Books, literature, poetry", "#1c1028", "#c898f0", "#9870c0"},
	{enums.Petsburg, "Petsburg", "Animals and people who love them", "#102018", "#60e888", "#40b868"},
	{enums.Purgatory, "Purgatory", "Horror, dark, gothic", "#180c30", "#b080ff", "#8858d0"},
	{enums.SiliconValley, "Silicon Valley", "Tech, programming, computers", "#102028", "#50d8c8", "#38a898"},
	{enums.SilverLake, "Silver Lake", "Music, bands, concerts", "#141c30", "#80b8f0", "#5888c0"},
	{enums.StratfordUponAvon, "Stratford-upon-Avon", "Writers and their writing", "#281010", "#ff8870", "#d06048"},
	{enums.Tokyo, "Tokyo", "Anime, manga, and the far east", "#280818", "#ff70a8", "#d04880"},
}

var lookup map[enums.DistrictSlug]District

func init() {
	lookup = make(map[enums.DistrictSlug]District, len(All))
	for _, entry := range All {
		lookup[entry.Slug] = entry
	}
}

func Find(slug enums.DistrictSlug) (District, bool) {
	entry, ok := lookup[slug]
	return entry, ok
}

func IsValid(slug string) bool {
	_, ok := lookup[enums.DistrictSlug(slug)]
	return ok
}