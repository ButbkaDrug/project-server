package reaper

import (
	"lp/pkg/comments"
	"lp/pkg/lyrics"
	"testing"
)

func TestAlternativeSection(t *testing.T) {
	OneThingRemains := lyrics.Lyrics{
		Name: "One Thing Remains",
		Sections: []lyrics.Section{
			{
				Name:     "Verse",
				Number:   1,
				Lang:     "",
				Comments: []comments.Comment{},
				Text: []lyrics.TextLine{
					{Offset: 0, Length: 0, Text: "Higher than the mountains that I face"},
					{Offset: 0, Length: 0, Text: "Stronger than the power of the grave"},
					{Offset: 0, Length: 0, Text: "Constant in the trial and the change"},
					{Offset: 0, Length: 0, Text: "This one thing remains"},
				},
			},
			{
				Name:     "Chorus",
				Number:   0,
				Lang:     "",
				Comments: []comments.Comment{},
				Text: []lyrics.TextLine{
					{Offset: 0, Length: 0, Text: "Your love never fails"},
					{Offset: 0, Length: 0, Text: "It never gives up"},
					{Offset: 0, Length: 0, Text: "It never runs out on me"},
					{Offset: 0, Length: 0, Text: "Your love never fails"},
					{Offset: 0, Length: 0, Text: "It never gives up"},
					{Offset: 0, Length: 0, Text: "It never runs out on me"},
					{Offset: 0, Length: 0, Text: "Your love never fails"},
					{Offset: 0, Length: 0, Text: "It never gives up"},
					{Offset: 0, Length: 0, Text: "It never runs out on me"},
				},
			},
			{
				Name:     "Verse",
				Number:   2,
				Lang:     "",
				Comments: []comments.Comment{},
				Text: []lyrics.TextLine{
					{Offset: 0, Length: 0, Text: "On and on and on and on it goes"},
					{Offset: 0, Length: 0, Text: "It overwhelms and satisfies my soul"},
					{Offset: 0, Length: 0, Text: "And I never ever have to be afraid"},
					{Offset: 0, Length: 0, Text: "This one thing remains"},
					{Offset: 0, Length: 0, Text: "This one thing remains"},
				},
			},
			{
				Name:     "Bridge",
				Number:   0,
				Lang:     "",
				Comments: []comments.Comment{},
				Text: []lyrics.TextLine{
					{Offset: 0, Length: 0, Text: "In death In life"},
					{Offset: 0, Length: 0, Text: "I'm confident and covered by"},
					{Offset: 0, Length: 0, Text: "The power of Your great love"},
					{Offset: 0, Length: 0, Text: "My debt is paid"},
					{Offset: 0, Length: 0, Text: "There's nothing that can separate"},
					{Offset: 0, Length: 0, Text: "My heart from Your great love"},
				},
			},
			{
				Name:     "Verse",
				Number:   1,
				Lang:     "R",
				Comments: []comments.Comment{},
				Text: []lyrics.TextLine{
					{Offset: 0, Length: 0, Text: "Выше всяких гор, что на пути."},
					{Offset: 0, Length: 0, Text: "Сильнее силы смерти, власти тьмы."},
					{Offset: 0, Length: 0, Text: "Скала среди бушующих морей"},
					{Offset: 0, Length: 0, Text: "Любовь сильней –"},
					{Offset: 0, Length: 0, Text: "раза"},
				},
			},
			{
				Name:     "Chorus",
				Number:   0,
				Lang:     "R",
				Comments: []comments.Comment{},
				Text: []lyrics.TextLine{
					{Offset: 0, Length: 0, Text: "Любовь Твоя"},
					{Offset: 0, Length: 0, Text: "Не подведёт"},
					{Offset: 0, Length: 0, Text: "И не предаст меня"},
					{Offset: 0, Length: 0, Text: "Любовь Твоя"},
					{Offset: 0, Length: 0, Text: "Не подведёт"},
					{Offset: 0, Length: 0, Text: "И не предаст меня"},
					{Offset: 0, Length: 0, Text: "Любовь Твоя"},
					{Offset: 0, Length: 0, Text: "Не подведёт"},
					{Offset: 0, Length: 0, Text: "И не предаст меня"},
				},
			},
			{
				Name:     "Verse",
				Number:   2,
				Lang:     "R",
				Comments: []comments.Comment{},
				Text: []lyrics.TextLine{
					{Offset: 0, Length: 0, Text: "Вновь, и вновь, и вновь, и вновь она"},
					{Offset: 0, Length: 0, Text: "Прольётся с неба, напоит меня"},
					{Offset: 0, Length: 0, Text: "Больше нечего теперь бояться мне"},
					{Offset: 0, Length: 0, Text: "Любовь сильней –"},
					{Offset: 0, Length: 0, Text: "раза"},
				},
			},
			{
				Name:     "Bridge",
				Number:   0,
				Lang:     "R",
				Comments: []comments.Comment{},
				Text: []lyrics.TextLine{
					{Offset: 0, Length: 0, Text: "И я всегда"},
					{Offset: 0, Length: 0, Text: "Уверен и наполнен силой"},
					{Offset: 0, Length: 0, Text: "Твоей любви"},
					{Offset: 0, Length: 0, Text: "Мой грех прощён,"},
					{Offset: 0, Length: 0, Text: "И ничто не может отделить"},
					{Offset: 0, Length: 0, Text: "Меня от Твоей любви."},
				},
			},
		},
	}

	tests := []Name{
		{"Verse", 1, "R"},
		{"Chorus", 0, "R"},
		{"Verse", 2, "R"},
		{"Bridge", 0, "R"},
		{"Verse", 1, ""},
		{"Chorus", 0, ""},
		{"Verse", 2, ""},
		{"Bridge", 0, ""},
	}

	for i, test := range tests {
		alt := OneThingRemains.GetAlternativeSection(OneThingRemains.Sections[i])
		if test.GetFullName() != alt.GetFullName() {
			t.Fatalf("test[%d] faild expected %s got %s", i, test.GetFullName(), OneThingRemains.Sections[i].GetFullName())
		}
	}
}

func TestRelativeSection(t *testing.T) {
	OneThingRemains := lyrics.Lyrics{
		Name: "One Thing Remains",
		Sections: []lyrics.Section{
			{
				Name:     "Verse",
				Number:   1,
				Lang:     "",
				Comments: []comments.Comment{},
				Text: []lyrics.TextLine{
					{Offset: 0, Length: 0, Text: "Higher than the mountains that I face"},
					{Offset: 0, Length: 0, Text: "Stronger than the power of the grave"},
					{Offset: 0, Length: 0, Text: "Constant in the trial and the change"},
					{Offset: 0, Length: 0, Text: "This one thing remains"},
				},
			},
			{
				Name:     "Chorus",
				Number:   0,
				Lang:     "",
				Comments: []comments.Comment{},
				Text: []lyrics.TextLine{
					{Offset: 0, Length: 0, Text: "Your love never fails"},
					{Offset: 0, Length: 0, Text: "It never gives up"},
					{Offset: 0, Length: 0, Text: "It never runs out on me"},
					{Offset: 0, Length: 0, Text: "Your love never fails"},
					{Offset: 0, Length: 0, Text: "It never gives up"},
					{Offset: 0, Length: 0, Text: "It never runs out on me"},
					{Offset: 0, Length: 0, Text: "Your love never fails"},
					{Offset: 0, Length: 0, Text: "It never gives up"},
					{Offset: 0, Length: 0, Text: "It never runs out on me"},
				},
			},
			{
				Name:     "Verse",
				Number:   2,
				Lang:     "",
				Comments: []comments.Comment{},
				Text: []lyrics.TextLine{
					{Offset: 0, Length: 0, Text: "On and on and on and on it goes"},
					{Offset: 0, Length: 0, Text: "It overwhelms and satisfies my soul"},
					{Offset: 0, Length: 0, Text: "And I never ever have to be afraid"},
					{Offset: 0, Length: 0, Text: "This one thing remains"},
					{Offset: 0, Length: 0, Text: "This one thing remains"},
				},
			},
			{
				Name:     "Bridge",
				Number:   0,
				Lang:     "",
				Comments: []comments.Comment{},
				Text: []lyrics.TextLine{
					{Offset: 0, Length: 0, Text: "In death In life"},
					{Offset: 0, Length: 0, Text: "I'm confident and covered by"},
					{Offset: 0, Length: 0, Text: "The power of Your great love"},
					{Offset: 0, Length: 0, Text: "My debt is paid"},
					{Offset: 0, Length: 0, Text: "There's nothing that can separate"},
					{Offset: 0, Length: 0, Text: "My heart from Your great love"},
				},
			},
			{
				Name:     "Verse",
				Number:   1,
				Lang:     "R",
				Comments: []comments.Comment{},
				Text: []lyrics.TextLine{
					{Offset: 0, Length: 0, Text: "Выше всяких гор, что на пути."},
					{Offset: 0, Length: 0, Text: "Сильнее силы смерти, власти тьмы."},
					{Offset: 0, Length: 0, Text: "Скала среди бушующих морей"},
					{Offset: 0, Length: 0, Text: "Любовь сильней –"},
					{Offset: 0, Length: 0, Text: "раза"},
				},
			},
			{
				Name:     "Chorus",
				Number:   0,
				Lang:     "R",
				Comments: []comments.Comment{},
				Text: []lyrics.TextLine{
					{Offset: 0, Length: 0, Text: "Любовь Твоя"},
					{Offset: 0, Length: 0, Text: "Не подведёт"},
					{Offset: 0, Length: 0, Text: "И не предаст меня"},
					{Offset: 0, Length: 0, Text: "Любовь Твоя"},
					{Offset: 0, Length: 0, Text: "Не подведёт"},
					{Offset: 0, Length: 0, Text: "И не предаст меня"},
					{Offset: 0, Length: 0, Text: "Любовь Твоя"},
					{Offset: 0, Length: 0, Text: "Не подведёт"},
					{Offset: 0, Length: 0, Text: "И не предаст меня"},
				},
			},
			{
				Name:     "Verse",
				Number:   2,
				Lang:     "R",
				Comments: []comments.Comment{},
				Text: []lyrics.TextLine{
					{Offset: 0, Length: 0, Text: "Вновь, и вновь, и вновь, и вновь она"},
					{Offset: 0, Length: 0, Text: "Прольётся с неба, напоит меня"},
					{Offset: 0, Length: 0, Text: "Больше нечего теперь бояться мне"},
					{Offset: 0, Length: 0, Text: "Любовь сильней –"},
					{Offset: 0, Length: 0, Text: "раза"},
				},
			},
			{
				Name:     "Bridge",
				Number:   0,
				Lang:     "R",
				Comments: []comments.Comment{},
				Text: []lyrics.TextLine{
					{Offset: 0, Length: 0, Text: "И я всегда"},
					{Offset: 0, Length: 0, Text: "Уверен и наполнен силой"},
					{Offset: 0, Length: 0, Text: "Твоей любви"},
					{Offset: 0, Length: 0, Text: "Мой грех прощён,"},
					{Offset: 0, Length: 0, Text: "И ничто не может отделить"},
					{Offset: 0, Length: 0, Text: "Меня от Твоей любви."},
				},
			},
		},
	}

	tests := []Name{
		{"Verse", 2, ""},
		{"Chorus", 0, ""},
		{"Verse", 1, ""},
		{"Bridge", 0, ""},
		{"Verse", 2, "R"},
		{"Chorus", 0, "R"},
		{"Verse", 1, "R"},
		{"Bridge", 0, "R"},
	}

	for i, test := range tests {
		alt := OneThingRemains.GetRelativeSection(OneThingRemains.Sections[i])
		if test.GetFullName() != alt.GetFullName() {
			t.Fatalf("test[%d] faild expected %s got %s", i, test.GetFullName(), OneThingRemains.Sections[i].GetFullName())
		}
	}
}

func TestRelativeAlternativeSection(t *testing.T) {
	OneThingRemains := lyrics.Lyrics{
		Name: "One Thing Remains",
		Sections: []lyrics.Section{
			{
				Name:     "Verse",
				Number:   1,
				Lang:     "",
				Comments: []comments.Comment{},
				Text: []lyrics.TextLine{
					{Offset: 0, Length: 0, Text: "Higher than the mountains that I face"},
					{Offset: 0, Length: 0, Text: "Stronger than the power of the grave"},
					{Offset: 0, Length: 0, Text: "Constant in the trial and the change"},
					{Offset: 0, Length: 0, Text: "This one thing remains"},
				},
			},
			{
				Name:     "Chorus",
				Number:   0,
				Lang:     "",
				Comments: []comments.Comment{},
				Text: []lyrics.TextLine{
					{Offset: 0, Length: 0, Text: "Your love never fails"},
					{Offset: 0, Length: 0, Text: "It never gives up"},
					{Offset: 0, Length: 0, Text: "It never runs out on me"},
					{Offset: 0, Length: 0, Text: "Your love never fails"},
					{Offset: 0, Length: 0, Text: "It never gives up"},
					{Offset: 0, Length: 0, Text: "It never runs out on me"},
					{Offset: 0, Length: 0, Text: "Your love never fails"},
					{Offset: 0, Length: 0, Text: "It never gives up"},
					{Offset: 0, Length: 0, Text: "It never runs out on me"},
				},
			},
			{
				Name:     "Verse",
				Number:   2,
				Lang:     "",
				Comments: []comments.Comment{},
				Text: []lyrics.TextLine{
					{Offset: 0, Length: 0, Text: "On and on and on and on it goes"},
					{Offset: 0, Length: 0, Text: "It overwhelms and satisfies my soul"},
					{Offset: 0, Length: 0, Text: "And I never ever have to be afraid"},
					{Offset: 0, Length: 0, Text: "This one thing remains"},
					{Offset: 0, Length: 0, Text: "This one thing remains"},
				},
			},
			{
				Name:     "Bridge",
				Number:   0,
				Lang:     "",
				Comments: []comments.Comment{},
				Text: []lyrics.TextLine{
					{Offset: 0, Length: 0, Text: "In death In life"},
					{Offset: 0, Length: 0, Text: "I'm confident and covered by"},
					{Offset: 0, Length: 0, Text: "The power of Your great love"},
					{Offset: 0, Length: 0, Text: "My debt is paid"},
					{Offset: 0, Length: 0, Text: "There's nothing that can separate"},
					{Offset: 0, Length: 0, Text: "My heart from Your great love"},
				},
			},
			{
				Name:     "Verse",
				Number:   1,
				Lang:     "R",
				Comments: []comments.Comment{},
				Text: []lyrics.TextLine{
					{Offset: 0, Length: 0, Text: "Выше всяких гор, что на пути."},
					{Offset: 0, Length: 0, Text: "Сильнее силы смерти, власти тьмы."},
					{Offset: 0, Length: 0, Text: "Скала среди бушующих морей"},
					{Offset: 0, Length: 0, Text: "Любовь сильней –"},
					{Offset: 0, Length: 0, Text: "раза"},
				},
			},
			{
				Name:     "Chorus",
				Number:   0,
				Lang:     "R",
				Comments: []comments.Comment{},
				Text: []lyrics.TextLine{
					{Offset: 0, Length: 0, Text: "Любовь Твоя"},
					{Offset: 0, Length: 0, Text: "Не подведёт"},
					{Offset: 0, Length: 0, Text: "И не предаст меня"},
					{Offset: 0, Length: 0, Text: "Любовь Твоя"},
					{Offset: 0, Length: 0, Text: "Не подведёт"},
					{Offset: 0, Length: 0, Text: "И не предаст меня"},
					{Offset: 0, Length: 0, Text: "Любовь Твоя"},
					{Offset: 0, Length: 0, Text: "Не подведёт"},
					{Offset: 0, Length: 0, Text: "И не предаст меня"},
				},
			},
			{
				Name:     "Verse",
				Number:   2,
				Lang:     "R",
				Comments: []comments.Comment{},
				Text: []lyrics.TextLine{
					{Offset: 0, Length: 0, Text: "Вновь, и вновь, и вновь, и вновь она"},
					{Offset: 0, Length: 0, Text: "Прольётся с неба, напоит меня"},
					{Offset: 0, Length: 0, Text: "Больше нечего теперь бояться мне"},
					{Offset: 0, Length: 0, Text: "Любовь сильней –"},
					{Offset: 0, Length: 0, Text: "раза"},
				},
			},
			{
				Name:     "Bridge",
				Number:   0,
				Lang:     "R",
				Comments: []comments.Comment{},
				Text: []lyrics.TextLine{
					{Offset: 0, Length: 0, Text: "И я всегда"},
					{Offset: 0, Length: 0, Text: "Уверен и наполнен силой"},
					{Offset: 0, Length: 0, Text: "Твоей любви"},
					{Offset: 0, Length: 0, Text: "Мой грех прощён,"},
					{Offset: 0, Length: 0, Text: "И ничто не может отделить"},
					{Offset: 0, Length: 0, Text: "Меня от Твоей любви."},
				},
			},
		},
	}

	tests := []Name{
		{"Verse", 2, "R"},
		{"Chorus", 0, ""},
		{"Verse", 1, "R"},
		{"Bridge", 0, ""},
		{"Verse", 2, ""},
		{"Chorus", 0, "R"},
		{"Verse", 1, ""},
		{"Bridge", 0, "R"},
	}

	for i, test := range tests {
		alt := OneThingRemains.GetRelativeAlternativeSection(OneThingRemains.Sections[i])
		if test.GetFullName() != alt.GetFullName() {
			t.Fatalf("test[%d] faild expected %s got %s", i, test.GetFullName(), OneThingRemains.Sections[i].GetFullName())
		}
	}
}

func TestAllRelativeSections(t *testing.T) {
	OneThingRemains := lyrics.Lyrics{
		Name: "One Thing Remains",
		Sections: []lyrics.Section{
			{
				Name:     "Verse",
				Number:   1,
				Lang:     "",
				Comments: []comments.Comment{},
				Text: []lyrics.TextLine{
					{Offset: 0, Length: 0, Text: "Higher than the mountains that I face"},
					{Offset: 0, Length: 0, Text: "Stronger than the power of the grave"},
					{Offset: 0, Length: 0, Text: "Constant in the trial and the change"},
					{Offset: 0, Length: 0, Text: "This one thing remains"},
				},
			},
			{
				Name:     "Chorus",
				Number:   0,
				Lang:     "",
				Comments: []comments.Comment{},
				Text: []lyrics.TextLine{
					{Offset: 0, Length: 0, Text: "Your love never fails"},
					{Offset: 0, Length: 0, Text: "It never gives up"},
					{Offset: 0, Length: 0, Text: "It never runs out on me"},
					{Offset: 0, Length: 0, Text: "Your love never fails"},
					{Offset: 0, Length: 0, Text: "It never gives up"},
					{Offset: 0, Length: 0, Text: "It never runs out on me"},
					{Offset: 0, Length: 0, Text: "Your love never fails"},
					{Offset: 0, Length: 0, Text: "It never gives up"},
					{Offset: 0, Length: 0, Text: "It never runs out on me"},
				},
			},
			{
				Name:     "Verse",
				Number:   2,
				Lang:     "",
				Comments: []comments.Comment{},
				Text: []lyrics.TextLine{
					{Offset: 0, Length: 0, Text: "On and on and on and on it goes"},
					{Offset: 0, Length: 0, Text: "It overwhelms and satisfies my soul"},
					{Offset: 0, Length: 0, Text: "And I never ever have to be afraid"},
					{Offset: 0, Length: 0, Text: "This one thing remains"},
					{Offset: 0, Length: 0, Text: "This one thing remains"},
				},
			},
			{
				Name:     "Bridge",
				Number:   0,
				Lang:     "",
				Comments: []comments.Comment{},
				Text: []lyrics.TextLine{
					{Offset: 0, Length: 0, Text: "In death In life"},
					{Offset: 0, Length: 0, Text: "I'm confident and covered by"},
					{Offset: 0, Length: 0, Text: "The power of Your great love"},
					{Offset: 0, Length: 0, Text: "My debt is paid"},
					{Offset: 0, Length: 0, Text: "There's nothing that can separate"},
					{Offset: 0, Length: 0, Text: "My heart from Your great love"},
				},
			},
			{
				Name:     "Verse",
				Number:   1,
				Lang:     "R",
				Comments: []comments.Comment{},
				Text: []lyrics.TextLine{
					{Offset: 0, Length: 0, Text: "Выше всяких гор, что на пути."},
					{Offset: 0, Length: 0, Text: "Сильнее силы смерти, власти тьмы."},
					{Offset: 0, Length: 0, Text: "Скала среди бушующих морей"},
					{Offset: 0, Length: 0, Text: "Любовь сильней –"},
					{Offset: 0, Length: 0, Text: "раза"},
				},
			},
			{
				Name:     "Chorus",
				Number:   0,
				Lang:     "R",
				Comments: []comments.Comment{},
				Text: []lyrics.TextLine{
					{Offset: 0, Length: 0, Text: "Любовь Твоя"},
					{Offset: 0, Length: 0, Text: "Не подведёт"},
					{Offset: 0, Length: 0, Text: "И не предаст меня"},
					{Offset: 0, Length: 0, Text: "Любовь Твоя"},
					{Offset: 0, Length: 0, Text: "Не подведёт"},
					{Offset: 0, Length: 0, Text: "И не предаст меня"},
					{Offset: 0, Length: 0, Text: "Любовь Твоя"},
					{Offset: 0, Length: 0, Text: "Не подведёт"},
					{Offset: 0, Length: 0, Text: "И не предаст меня"},
				},
			},
			{
				Name:     "Verse",
				Number:   2,
				Lang:     "R",
				Comments: []comments.Comment{},
				Text: []lyrics.TextLine{
					{Offset: 0, Length: 0, Text: "Вновь, и вновь, и вновь, и вновь она"},
					{Offset: 0, Length: 0, Text: "Прольётся с неба, напоит меня"},
					{Offset: 0, Length: 0, Text: "Больше нечего теперь бояться мне"},
					{Offset: 0, Length: 0, Text: "Любовь сильней –"},
					{Offset: 0, Length: 0, Text: "раза"},
				},
			},
			{
				Name:     "Bridge",
				Number:   0,
				Lang:     "R",
				Comments: []comments.Comment{},
				Text: []lyrics.TextLine{
					{Offset: 0, Length: 0, Text: "И я всегда"},
					{Offset: 0, Length: 0, Text: "Уверен и наполнен силой"},
					{Offset: 0, Length: 0, Text: "Твоей любви"},
					{Offset: 0, Length: 0, Text: "Мой грех прощён,"},
					{Offset: 0, Length: 0, Text: "И ничто не может отделить"},
					{Offset: 0, Length: 0, Text: "Меня от Твоей любви."},
				},
			},
		},
	}

	tests := []struct {
		Name   Name
		Number int
	}{

		{Name{"Verse", 1, ""}, 4},
		{Name{"Chorus", 0, ""}, 2},
		{Name{"Verse", 2, ""}, 4},
		{Name{"Bridge", 0, ""}, 2},
		{Name{"Verse", 1, "R"}, 4},
		{Name{"Chorus", 0, "R"}, 2},
		{Name{"Verse", 2, "R"}, 4},
		{Name{"Bridge", 0, "R"}, 2},
	}

	for i, test := range tests {
		alt := OneThingRemains.GetAllRelativeSections(OneThingRemains.Sections[i])
		if test.Number != len(alt) {
			t.Fatalf("test[%d] faild expected %v got %v %v", i, test, alt, len(alt))
		}
	}
}
