package reaper

import (
	"lp/pkg/comments"
	"lp/pkg/lyrics"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewProjectWithConfig(t *testing.T) {

	//first let's try empty config
	config := ProjectConfig{}
	assert.NotNil(t, config)

	reaper := NewProjectWithConfig(&config)

	if reaper == nil {
		t.Fatalf("NewProjectWithConfig failed. expected %T got %T", &Project{}, reaper)
	}

	if reaper.Path != config.Path {
		t.Fatalf("NewProjectWithConfig failed. expected Path=%v got %v", config.Path, reaper.Path)
	}

	if reaper.TrackID != config.TrackID {
		t.Fatalf("NewProjectWithConfig failed. expected TrackID=%v got %v", config.TrackID, reaper.TrackID)
	}

	if reaper.logger == nil {
		t.Fatalf("NewProjectWithConfig failed. expected logger got nil")
	}

}

func TestUpdateCursor(t *testing.T) {
	proj := NewProject()

	tests := []float64{
		0,
		5.10637188208617,
		5.10637188208617,
		15.31913832199547,
		14.04256235827664,
		15.31913832199547,
		35.74467120181406,
		17.92667535099146,
		19.14893424036281,
		24.25532879818594,
		24.44544796642061,
		29.36170068027211,
		30.96424952959907,
		33.82979591836735,
		34.46808510638298,
		35.74467120181406,
		56.17020408163265,
		37.48302214502823,
		39.5744671201814,
		43.18697351280938,
		44.00182370820669,
		44.68086167800453,
		49.45977324263038,
		50.52059632363583,
		54.25532879818594,
		54.89362811791383,
		56.17020408163265,
		66.38297052154195,
		57.0393978868143,
		60,
		63.14013605442177,
		65.10637188208617,
		65.74467120181406,
		66.38297052154195,
		97.02126984126984,
		69.26212187002461,
		70.0769720654219,
		70.85106575963719,
		74.98462585034014,
		76.59574468085107,
		81.06383219954648,
		83.11451729628021,
		83.92936749167752,
		86.17020408163265,
		90.27834467120181,
		90.448169054856,
		95.4014739229025,
		95.72505427702997,
		95.74467120181406,
		97.02126984126984,
		107.23403628117914,
		100.85106575963719,
		103.48574323346358,
		105.95743764172336,
		107.23403628117914,
		127.65956916099773,
		110.00451584889274,
		111.06383219954648,
		115.2489506440874,
		116.17020408163265,
		120.85510204081632,
		121.78911564625849,
		122.22723983210305,
		126.09113378684808,
		126.38297052154195,
		127.65956916099773,
		137.87233560090704,
		128.74604139528154,
		131.48936507936509,
		135.26481401071067,
		136.93380952380951,
		137.23403628117913,
		137.87233560090704,
		168.51063492063491,
		141.78361557388911,
		142.34043083900227,
		146.80850340136055,
		148.30238818931826,
		152.55319727891157,
		154.28310898827615,
		157.65956916099773,
		160.96740483427413,
		161.33996236792589,
		162.12766439909296,
		166.59573696145125,
		166.91376417233559,
		167.85876393110436,
		168.51063492063491,
		178.72340136054422,
		172.34043083900227,
		174.80911854103343,
		175.19238674193079,
		177.44680272108843,
		178.72340136054422,
		188.93616780045352,
		181.71118830510929,
		187.41511072514112,
		188.93616780045352,
		209.36170068027212,
		188.93616780045352,
		194.04256235827665,
		194.74876248371689,
		198.82727891156463,
		201.26753509914604,
		202.97873015873017,
		206.97148646692716,
		207.44680272108843,
		209.36170068027212,
		229.7872335600907,
		209.36170068027212,
		212.67540888695905,
		213.08140107106672,
		214.28673469387755,
		219.20308390022674,
		220.00906064553482,
		223.23133786848072,
		226.52783326096397,
		227.87233560090704,
		229.7872335600907,
		250.2127664399093,
		241.1951078303662,
		247.71390939354464,
		250.2127664399093,
		270.63829931972788,
		253.82205818497613,
		259.12178318135767,
		264.82570560138953,
		267.27025618758142,
		270.63829931972788,
		301.27659863945576,
		273.5575336517586,
		279.83372412794904,
		284.97617600231581,
		290.90085395860473,
		301.27659863945576,
		311.48936507936509,
		1304.6808510638298,
		1304.6808510638298,
		1320.8204371110146,
		1320.8204371110146,
		1320.8204371110146,
		1342.3398755246781,
		1326.2002894774932,
		1334.9425676653639,
		1342.3398755246778,
		1342.3398755246781,
		1363.8593139383413,
		1347.7197278911563,
		1356.4620060790273,
		1361.1693732812275,
		1361.8418584455057,
		1363.8593139383413,
		1385.3787523520045,
		1367.2217108119844,
		1373.2740483427413,
		1378.6539007092199,
		1384.0337820234477,
		1384.7062671877261,
		1385.3787523520045,
		1406.8981618179187,
		1394.7089303806629,
		1400.0628745115068,
		1405.5532204371111,
		1406.2257056013896,
		1406.8981618179187,
		1431.1075408886959,
		1416.312925170068,
		1421.3216963381094,
		1429.7625705601388,
		1431.1075408886957,
		1431.1075408886959,
		1452.6269793023594,
		1436.4873932551745,
		1445.2296714430454,
		1449.9370386452451,
		1452.6269793023591,
		1452.6269793023594,
		1474.1464177160226,
		1458.0068316688378,
		1466.9588941959764,
		1470.7839918946299,
		1474.1464177160226,
		1495.6658561296858,
		1495.6658561296858,
		1517.1852945433493,
		1517.1852945433493,
		1541.3946446663772,
		1541.3946446663772,
		1562.9140830800407,
		1562.9140830800407,
		1587.1234621508179,
		1587.1234621508179,
		1597.8831668837747,
		1596.5381965552181,
		1597.8831668837747,
		1619.4026052974382,
		1605.9529599073671,
		1607.2979302359242,
		1616.7126935880738,
		1619.4026052974382,
		1643.6119843682154,
		1643.6119843682154,
		1665.1314227818789,
		1665.1314227818789,
		1686.6508611955421,
		1686.6508611955421,
		1708.1702996092054,
		1708.1702996092054,
		1732.3796497322335,
		1732.3796497322335,
		1743.1393834129399,
		1742.4668982486612,
		1743.1393834129399,
		1753.8990881458967,
		1752.5541178173398,
		2614.4680851063831,
		2614.4680851063831,
		2623.4793457808655,
		2623.4793457808655,
		2641.5018671298308,
		2640.9427925531918,
		2641.5018671298308,
		2677.5469387755102,
		2644.881090425532,
		2649.9499468085105,
		2653.8923670212762,
		2658.9611968085105,
		2662.9036170212762,
		2667.9724734042552,
		2671.9148936170213,
		2675.2941223404259,
		2675.857313829787,
		2677.5469387755102,
		2713.5919814734407,
		2685.4317819148937,
		2694.4430585106384,
		2703.4543085106384,
		2712.4655851063826,
		2713.0287765957451,
		2713.5919814734407,
		2749.6370531191201,
		2716.9712234042549,
		2722.0400531914893,
		2725.982473404255,
		2731.0513031914893,
		2734.99375,
		2740.062579787234,
		2747.3842287234047,
		2749.6370531191201,
		2785.6820958170501,
		2785.6820958170501,
		2821.7271674627295,
		2821.7271674627295,
		2839.7496888116948,
		2822.51125,
		2825.6695744680851,
		2831.3016223404252,
		2834.6808510638298,
		2839.7496888116948,
		2857.7722101606601,
		2857.7722101606601,
		2875.7947315096249,
		2875.7947315096249,
		2893.8172818063395,
		2893.8172818063395,
		2929.8623245042695,
		2929.8623245042695,
		2947.8848458532348,
		2947.8848458532348,
		2965.9073961499494,
		2965.9073961499494,
		2983.9299174989142,
		3921.7021276595747,
		3921.7021276595747,
		3937.7337964973226,
		3937.7337964973226,
		3959.1093645969027,
		3937.7538862353454,
		3942.6323635837311,
		3948.8668982486615,
		3952.8748009842234,
		3957.7733680706319,
		3959.1093645969027,
		3964.4532349109859,
		3963.5625995078881,
		3964.4532349109859,
		3985.8288030105664,
		3968.9064698219713,
		3974.2503690838039,
		3984.492835432045,
		3985.8288030105664,
		4007.2043421623971,
		3986.2741351859895,
		3990.7273700969749,
		3996.9619047619049,
		4000.9698074974667,
		4005.8683745838762,
		4007.2043421623971,
		4012.5482414242297,
		4012.5482414242297,
		4033.9238095238097,
		4033.033145172963,
		4033.9238095238097,
		4055.2993486756409,
		4038.3770444347947,
		4043.7209436966277,
		4053.96338109712,
		4054.4087132725431,
		4055.2993486756409,
		4076.6749167752209,
		4060.1979157620499,
		4065.0964828484589,
		4070.8857142857141,
		4075.3389491967,
		4076.6749167752209,
		4098.0504559270521,
		4098.0504559270521,
		4103.3943551888842,
		4103.3943551888842,
		4124.7699232884652,
		4124.7699232884652,
		4146.1454624402959,
		4146.1454624402959,
		4167.5210305398759,
		4167.5210305398759,
		4188.896598639456,
		4188.896598639456,
		4210.2721377912876,
		4210.2721377912876,
		4231.6477058908677,
		4231.6477058908677,
		4253.0232450426984,
		4253.0232450426984,
		4274.3988131422784,
		4274.3988131422784,
		4279.7427124041114,
		4279.7427124041114,
		4285.0865827181942,
		4285.0865827181942,
		4287.7585468229854,
		4287.7585468229854,
		4309.1340859748161,
		4309.1340859748161,
		4330.5096251266468,
		5228.9361702127662,
		6536.1702127659573,
		7843.4042553191493,
		15686.808510638299,
		17647.659574468085,
	}

	for i, test := range tests {
		proj.UpdateCursor(test)

		if proj.Cursor != test {
			t.Fatalf("UpdateCursor() test[%d] failed. expected %v got %v", i, test, proj.Cursor)

		}
	}
}

func GetCurrentLineTest(t *testing.T) {
	proj := NewProject()

	proj.Lines = []lyrics.LyricLine{
		{ID: 94, Position: 1456.9567073856742, Length: 5.6536462652, Text: "Имя святое его"},
		{ID: 95, Position: 1462.6103536508742, Length: 4.2773882175, Text: "Пой Ему, как никогда"},
		{ID: 96, Position: 1466.8877418683742, Length: 3.168300811, Text: "Пой душа"},
		{ID: 97, Position: 1470.056042679374, Length: 9.2454922925, Text: "Благослови Господа"},
		{ID: 98, Position: 1476.8187885638074, Length: 1503.3015385638075, Text: ""},
		{ID: 99, Position: 1502.8252103195075, Length: 2.7106906812, Text: "And on that day"},
		{ID: 100, Position: 1505.5359010007076, Length: 4.3863272183, Text: "When my strength is failing"},
		{ID: 101, Position: 1509.9222282190076, Length: 6.1750190239, Text: "The end draws near and my time has come"},
		{ID: 102, Position: 1516.0972472429075, Length: 2.9284292519, Text: "Still my soul will"},
		{ID: 103, Position: 1519.0256764948076, Length: 5.7931034483, Text: "Sing Your praise unending"},
		{ID: 104, Position: 1524.8187799431075, Length: 2.0606200734, Text: "Ten thousand years"},
		{ID: 105, Position: 1526.8794000165076, Length: 5.3876557887, Text: "And then forevermore"},
		{ID: 106, Position: 1535.4706205451407, Length: 4.2447211796, Text: "Bless the Lord O my soul"},
		{ID: 107, Position: 1539.7153417247407, Length: 3.3103448276, Text: "O my soul"},
		{ID: 108, Position: 1543.0256865523409, Length: 5.6536462652, Text: "Worship His Holy name"},
		{ID: 109, Position: 1548.6793328175409, Length: 4.2773882175, Text: "Sing like never before"},
		{ID: 110, Position: 1552.9567210350408, Length: 3.168300811, Text: "O my soul"},
		{ID: 111, Position: 1556.1250218460407, Length: 9.2454922925, Text: "I'll worship Your Holy name"},
		{ID: 112, Position: 1561.9533705451408, Length: 4.2447211796, Text: "Славь душа Господа"},
		{ID: 113, Position: 1566.1980917247408, Length: 3.3103448276, Text: "Cлавь душа"},
		{ID: 114, Position: 1569.508436552341, Length: 5.6536462652, Text: "Имя святое его"},
		{ID: 115, Position: 1575.162082817541, Length: 4.2773882175, Text: "Пой Ему, как никогда"},
		{ID: 116, Position: 1579.439471035041, Length: 3.168300811, Text: "Пой душа"},
		{ID: 117, Position: 1582.6077718460408, Length: 9.2454922925, Text: "Благослови Господа"},
		{ID: 118, Position: 2614.468085106346, Length: 2627.026225695915, Text: ""},
		{ID: 119, Position: 2627.041973895115, Length: 3.8214611031, Text: "Who am I that the highest King"},
		{ID: 120, Position: 2630.863434998215, Length: 4.8837209302, Text: "Would welcome me"},
		{ID: 121, Position: 2635.747155928415, Length: 3.1395348838, Text: "I was lost but He brought me in"},
		{ID: 122, Position: 2638.886690812215, Length: 3.8372093023, Text: "Oh His love for me"},
		{ID: 123, Position: 2643.770420707253, Length: 2647.9564524532852, Text: ""},
		{ID: 124, Position: 2647.258778034685, Length: 4.1860465116, Text: "Who the Son sets free"},
		{ID: 125, Position: 2651.444824546285, Length: 4.1860465117, Text: "Oh is free indeed"},
		{ID: 126, Position: 2655.6308710579847, Length: 8.0232558139, Text: "I'm a child of God yes I am"},
		{ID: 127, Position: 2665.0494846739225, Length: 3.4883720931, Text: "Я свободен искуплен Им"},
		{ID: 128, Position: 2668.5378567670223, Length: 4.8837209302, Text: "Он ми — ло — стив"},
		{ID: 129, Position: 2673.4215776972223, Length: 3.1395348837, Text: "Был я в рабстве греха, но Ты"},
		{ID: 130, Position: 2676.5611125809223, Length: 3.8372093023, Text: "Жизнью за — пла - тил"},
		{ID: 131, Position: 2681.4448198002233, Length: 2685.630874221992, Text: ""},
		{ID: 132, Position: 2684.933199803392, Length: 4.1860465116, Text: "Дал свободу мне"},
		{ID: 133, Position: 2689.1192463149923, Length: 4.1860465117, Text: "Иисус навек"},
		{ID: 134, Position: 2693.305292826692, Length: 8.0232558139, Text: "Я дитя Ца - ря знаю я"},
		{ID: 135, Position: 2701.67739481473, Length: 4.1860465116, Text: "В Доме есть всегда"},
		{ID: 136, Position: 2705.8634413263303, Length: 4.1860465116, Text: "Место для меня"},
		{ID: 137, Position: 2710.04948783793, Length: 8.023255814, Text: "Я дитя Ца - ря знаю я"},
		{ID: 138, Position: 2718.4215671503307, Length: 4.5348837209, Text: "Богом избран не оставлен"},
		{ID: 139, Position: 2722.9564508712306, Length: 3.8372093023, Text: "Я знаю, кто я в Тебе"},
		{ID: 140, Position: 2726.7936601735305, Length: 4.534883721, Text: "Твоим другом был я назван"},
		{ID: 141, Position: 2731.328543894531, Length: 3.488372093, Text: "Я знаю, кто я в Тебе"},
		{ID: 142, Position: 2735.1657621616687, Length: 4.5348837209, Text: "Богом избран не оставлен"},
		{ID: 143, Position: 2739.7006458825686, Length: 3.8372093023, Text: "Я знаю, кто я в Тебе"},
		{ID: 144, Position: 2743.5378551848685, Length: 4.534883721, Text: "Твоим другом был я назван"},
		{ID: 145, Position: 2748.072738905869, Length: 3.488372093, Text: "Я знаю, кто я в Тебе"},
		{ID: 146, Position: 2752.6076089158696, Length: 2756.7936633376385, Text: ""},
		{ID: 147, Position: 2756.0959889190385, Length: 4.1860465116, Text: "Дал свободу мне"},
		{ID: 148, Position: 2760.2820354306386, Length: 4.1860465117, Text: "Иисус навек"},
		{ID: 149, Position: 2764.4680819423384, Length: 8.0232558139, Text: "Я дитя Ца - ря знаю я"},
		{ID: 150, Position: 2772.8401839303765, Length: 4.1860465116, Text: "В Доме есть всегда"},
		{ID: 151, Position: 2777.0262304419766, Length: 4.1860465116, Text: "Место для меня"},
		{ID: 152, Position: 2781.2122769535763, Length: 8.023255814, Text: "Я дитя Ца - ря знаю я"},
		{ID: 153, Position: 2789.584356265977, Length: 4.1860465116, Text: "В Доме есть всегда"},
		{ID: 154, Position: 2793.770402777577, Length: 4.1860465116, Text: "Место для меня"},
		{ID: 155, Position: 2797.956449289177, Length: 8.023255814, Text: "Я дитя Ца - ря знаю я"},
		{ID: 156, Position: 2807.026225695915, Length: 2823.770420707253, Text: ""},
		{ID: 157, Position: 2823.072746288653, Length: 4.5348837209, Text: "I am chosen not forsaken"},
		{ID: 158, Position: 2827.607630009553, Length: 3.8372093023, Text: "I am who You say I am"},
		{ID: 159, Position: 2831.444839311853, Length: 4.534883721, Text: "You are for me not against me"},
		{ID: 160, Position: 2835.979723032853, Length: 3.488372093, Text: "I am who You say I am"},
		{ID: 161, Position: 2839.816918624254, Length: 4.5348837209, Text: "I am chosen not forsaken"},
		{ID: 162, Position: 2844.351802345154, Length: 3.8372093023, Text: "I am who You say I am"},
		{ID: 163, Position: 2848.1890116474538, Length: 4.534883721, Text: "You are for me not against me"},
		{ID: 164, Position: 2852.723895368454, Length: 3.488372093, Text: "I am who You say I am"},
		{ID: 165, Position: 39475445, Length: 2.8437148543, Text: "Иисус, иисус,"},
		{ID: 255, Position: 4226.167454329745, Length: 3.6427716322, Text: "Приводишь тьму ты в трепет!"},
		{ID: 256, Position: 4229.810225961945, Length: 2.8378378379, Text: "Иисус, иисус"},
		{ID: 257, Position: 4226.441253191024, Length: 3.3689929952, Text: "Иисус, иисус,"},
		{ID: 258, Position: 4229.810246186224, Length: 3.0910167098, Text: "Приводишь тьму ты в трепет!"},
		{ID: 259, Position: 4232.901262896024, Length: 3.3954697767, Text: "Иисус, иисус, и страх молчит!"},
		{ID: 260, Position: 4236.296732672724, Length: 2.8437148543, Text: "Иисус, иисус,"},
		{ID: 261, Position: 4239.140447527024, Length: 3.6427716322, Text: "Приводишь тьму ты в трепет!"},
		{ID: 262, Position: 4242.783219159224, Length: 2.8378378379, Text: "Иисус, иисус"},
		{ID: 263, Position: 4239.539973464565, Length: 4258.999429246878, Text: ""},
		{ID: 264, Position: 4258.999429246878, Length: 4.4594594594, Text: "Мир, принеси твой мир"},
		{ID: 265, Position: 4263.458888706278, Length: 2.4324324325, Text: "Бури успокой"},
		{ID: 266, Position: 4265.891321138778, Length: 5.6756756756, Text: "Повели им уйти!"},
		{ID: 267, Position: 4271.566996814378, Length: 4.8648648649, Text: "Бог, слушают тебя"},
		{ID: 268, Position: 4276.431861679278, Length: 2.4324324324, Text: "И бездны и моря"},
	}

	tests := []struct {
		Position   float64
		ExpectedID int
	}{
		{2663.65412687, 126},
		{2664.65412687, 126},
		{2665.03210000, 126},
		{2665.0494846739225, 127},
		//that what i get in reality {ID:123 Position:2643.770420707253 Length:2647.9564524532852 Text:}
		{2664.26025390625, 126},
	}

	for i, line := range proj.Lines {
		proj.Update(line.Position)

		if line.ID != proj.Line.ID {
			t.Fatalf("UpdateLine() test[%d] failed. expected %v got %v", i, line, proj.Line)
		}
	}

	for i, line := range proj.Lines {
		proj.Update(line.Position + line.Length)

		if line.ID == proj.Line.ID {
			t.Fatalf("UpdateLine() test[%d] failed. not expected %v but got %v", i, line, proj.Line)
		}
	}

	for i, test := range tests {
		proj.Update(test.Position)

		if proj.Line.ID != test.ExpectedID {
			t.Fatalf("UpdateLine() test[%d] failed. expected ID=%v got %v", i, test.ExpectedID, proj.Line.ID)
		}
	}
}

func TestParseSectionName(t *testing.T) {

	data := []string{
		"Verse",
		"Verse R",
		"Verse 1",
		"Verse 1 R",
		"Pre Chorus",
		"Pre Chorus R",
		"Pre Chorus 1",
		"Pre Chorus 2 R",
	}

	tests := []struct {
		Name   string
		Lang   string
		Number int
	}{
		{Name: "Verse"},
		{Name: "Verse", Lang: "R"},
		{Name: "Verse", Number: 1},
		{Name: "Verse", Lang: "R", Number: 1},
		{Name: "Pre Chorus"},
		{Name: "Pre Chorus", Lang: "R"},
		{Name: "Pre Chorus", Number: 1},
		{Name: "Pre Chorus", Lang: "R", Number: 2},
	}

	for i, test := range tests {

		n := parseSectionName(data[i])

		if n.Name != test.Name {
			t.Fatalf("test[%d] failed expected name %v got %v", i, test.Name, n.Name)
		}

		if n.Lang != test.Lang {
			t.Fatalf("test[%d] failed expected lang %v got %v", i, test.Lang, n.Lang)
		}

		if n.Number != test.Number {
			t.Fatalf("test[%d] failed expected number %v got %v", i, test.Number, n.Number)
		}

		if n.GetFullName() != data[i] {
			t.Fatalf("test[%d] failed expected %v got %v", i, data[i], n.GetFullName())
		}

	}
}

func TestCalculateOffsets(t *testing.T) {
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
					{
						Offset: 0,
						Length: 0,
						Text:   "Your love never fails"},
					{
						Offset: 0,
						Length: 0,
						Text:   "It never gives up"},
					{
						Offset: 0,
						Length: 0,
						Text:   "It never runs out on me"},
					{
						Offset: 0,
						Length: 0,
						Text:   "Your love never fails"},
					{
						Offset: 0,
						Length: 0,
						Text:   "It never gives up"},
					{
						Offset: 0,
						Length: 0,
						Text:   "It never runs out on me"},
					{
						Offset: 0,
						Length: 0,
						Text:   "Your love never fails"},
					{
						Offset: 0,
						Length: 0,
						Text:   "It never gives up"},
					{
						Offset: 0,
						Length: 0,
						Text:   "It never runs out on me"}}},
			{
				Name:     "Verse",
				Number:   2,
				Lang:     "",
				Comments: []comments.Comment{},
				Text: []lyrics.TextLine{
					{
						Offset: 0,
						Length: 0,
						Text:   "On and on and on and on it goes"},
					{
						Offset: 0,
						Length: 0,
						Text:   "It overwhelms and satisfies my soul"},
					{
						Offset: 0,
						Length: 0,
						Text:   "And I never ever have to be afraid"},
					{
						Offset: 0,
						Length: 0,
						Text:   "This one thing remains"},
					{
						Offset: 0,
						Length: 0,
						Text:   "This one thing remains"}}},
			{
				Name:     "Bridge",
				Number:   0,
				Lang:     "",
				Comments: []comments.Comment{},
				Text: []lyrics.TextLine{
					{
						Offset: 0,
						Length: 0,
						Text:   "In death In life"},
					{
						Offset: 0,
						Length: 0,
						Text:   "I'm confident and covered by"},
					{
						Offset: 0,
						Length: 0,
						Text:   "The power of Your great love"},
					{
						Offset: 0,
						Length: 0,
						Text:   "My debt is paid"},
					{
						Offset: 0,
						Length: 0,
						Text:   "There's nothing that can separate"},
					{
						Offset: 0,
						Length: 0,
						Text:   "My heart from Your great love"}}},
			{
				Name:     "Verse",
				Number:   1,
				Lang:     "R",
				Comments: []comments.Comment{},
				Text: []lyrics.TextLine{
					{
						Offset: 0,
						Length: 0,
						Text:   "Выше всяких гор, что на пути."},
					{
						Offset: 0,
						Length: 0,
						Text:   "Сильнее силы смерти, власти тьмы."},
					{
						Offset: 0,
						Length: 0,
						Text:   "Скала среди бушующих морей"},
					{
						Offset: 0,
						Length: 0,
						Text:   "Любовь сильней –"},
					{
						Offset: 0,
						Length: 0,
						Text:   "раза"}}},
			{
				Name:     "Chorus",
				Number:   0,
				Lang:     "R",
				Comments: []comments.Comment{},
				Text: []lyrics.TextLine{
					{
						Offset: 0,
						Length: 0,
						Text:   "Любовь Твоя"},
					{
						Offset: 0,
						Length: 0,
						Text:   "Не подведёт"},
					{
						Offset: 0,
						Length: 0,
						Text:   "И не предаст меня"},
					{
						Offset: 0,
						Length: 0,
						Text:   "Любовь Твоя"},
					{
						Offset: 0,
						Length: 0,
						Text:   "Не подведёт"},
					{
						Offset: 0,
						Length: 0,
						Text:   "И не предаст меня"},
					{
						Offset: 0,
						Length: 0,
						Text:   "Любовь Твоя"},
					{
						Offset: 0,
						Length: 0,
						Text:   "Не подведёт"},
					{
						Offset: 0,
						Length: 0,
						Text:   "И не предаст меня"}}},
			{
				Name:     "Verse",
				Number:   2,
				Lang:     "R",
				Comments: []comments.Comment{},
				Text: []lyrics.TextLine{
					{
						Offset: 0,
						Length: 0,
						Text:   "Вновь, и вновь, и вновь, и вновь она"},
					{
						Offset: 0,
						Length: 0,
						Text:   "Прольётся с неба, напоит меня"},
					{
						Offset: 0,
						Length: 0,
						Text:   "Больше нечего теперь бояться мне"},
					{
						Offset: 0,
						Length: 0,
						Text:   "Любовь сильней –"},
					{
						Offset: 0,
						Length: 0,
						Text:   "раза"}}},
			{
				Name:     "Bridge",
				Number:   0,
				Lang:     "R",
				Comments: []comments.Comment{},
				Text: []lyrics.TextLine{
					{
						Offset: 0,
						Length: 0,
						Text:   "И я всегда"},
					{
						Offset: 0,
						Length: 0,
						Text:   "Уверен и наполнен силой"},
					{
						Offset: 0,
						Length: 0,
						Text:   "Твоей любви"},
					{
						Offset: 0,
						Length: 0,
						Text:   "Мой грех прощён,"},
					{
						Offset: 0,
						Length: 0,
						Text:   "И ничто не может отделить"},
					{
						Offset: 0,
						Length: 0,
						Text:   "Меня от Твоей любви."}}}}}

	BlessTheLord := lyrics.Lyrics{
		Name: "10000 Reasons",
		Sections: []lyrics.Section{
			{
				Name:     "Chorus",
				Number:   0,
				Lang:     "",
				Comments: []comments.Comment{},
				Text: []lyrics.TextLine{
					{
						Offset: -0.934376352,
						Length: 4.2447211796,
						Text:   "Bless the Lord O my soul"},
					{
						Offset: 3.3103448276,
						Length: 3.3103448276,
						Text:   "O my soul"},
					{
						Offset: 6.6206896552,
						Length: 5.6536462652,
						Text:   "Worship His Holy name"},
					{
						Offset: 12.2743359204,
						Length: 4.2773882175,
						Text:   "Sing like never before"},
					{
						Offset: 16.5517241379,
						Length: 3.168300811,
						Text:   "O my soul"},
					{
						Offset: 19.7200249489,
						Length: 9.2454922925,
						Text:   "I'll worship Your Holy name"}}},
			{
				Name:     "Verse",
				Number:   1,
				Lang:     "",
				Comments: []comments.Comment{},
				Text: []lyrics.TextLine{
					{
						Offset: -0.4763282443,
						Length: 2.7106906812,
						Text:   "The sun comes up"},
					{
						Offset: 2.2343624369,
						Length: 4.3863272183,
						Text:   "It's a new day dawning"},
					{
						Offset: 6.6206896552,
						Length: 6.1750190239,
						Text:   "It's time to sing Your song again"},
					{
						Offset: 12.7957086791,
						Length: 2.9284292519,
						Text:   "Whatever may pass"},
					{
						Offset: 15.724137931,
						Length: 5.7931034483,
						Text:   "And whatever lies before me"},
					{
						Offset: 21.5172413793,
						Length: 2.0606200734,
						Text:   "Let me be singing"},
					{
						Offset: 23.5778614527,
						Length: 5.3876557887,
						Text:   "When the evening comes"}}},
			{
				Name:     "Verse",
				Number:   2,
				Lang:     "",
				Comments: []comments.Comment{},
				Text: []lyrics.TextLine{
					{
						Offset: -0.4763282443,
						Length: 2.7106906812,
						Text:   "You're rich in love"},
					{
						Offset: 2.2343624369,
						Length: 4.3863272183,
						Text:   "And You're slow to anger"},
					{
						Offset: 6.6206896552,
						Length: 6.1750190239,
						Text:   "Your name is great and Your heart is kind"},
					{
						Offset: 12.7957086791,
						Length: 2.9284292519,
						Text:   "For all Your goodness"},
					{
						Offset: 15.724137931,
						Length: 5.7931034483,
						Text:   "I will keep on singing"},
					{
						Offset: 21.5172413793,
						Length: 2.0606200734,
						Text:   "Ten thousand reasons"},
					{
						Offset: 23.5778614527,
						Length: 5.3876557887,
						Text:   "For my heart to find"}}},
			{
				Name:     "Verse",
				Number:   3,
				Lang:     "",
				Comments: []comments.Comment{},
				Text: []lyrics.TextLine{
					{
						Offset: -0.4763282443,
						Length: 2.7106906812,
						Text:   "And on that day"},
					{
						Offset: 2.2343624369,
						Length: 4.3863272183,
						Text:   "When my strength is failing"},
					{
						Offset: 6.6206896552,
						Length: 6.1750190239,
						Text:   "The end draws near and my time has come"},
					{
						Offset: 12.7957086791,
						Length: 2.9284292519,
						Text:   "Still my soul will"},
					{
						Offset: 15.724137931,
						Length: 5.7931034483,
						Text:   "Sing Your praise unending"},
					{
						Offset: 21.5172413793,
						Length: 2.0606200734,
						Text:   "Ten thousand years"},
					{
						Offset: 23.5778614527,
						Length: 5.3876557887,
						Text:   "And then forevermore"}}},
			{
				Name:     "Tag",
				Number:   0,
				Lang:     "",
				Comments: []comments.Comment{},
				Text: []lyrics.TextLine{
					{
						Offset: 0,
						Length: 0,
						Text:   "I'll worship Your Holy name"}}},
			{
				Name:     "Verse",
				Number:   1,
				Lang:     "R",
				Comments: []comments.Comment{},
				Text: []lyrics.TextLine{
					{
						Offset: -0.4763282443,
						Length: 2.7106906812,
						Text:   "Солнце встаёт"},
					{
						Offset: 2.2343624369,
						Length: 4.3863272183,
						Text:   "Новый день приходит"},
					{
						Offset: 6.6206896552,
						Length: 6.1750190239,
						Text:   "Время опять хвалу воздать"},
					{
						Offset: 12.7957086791,
						Length: 2.9284292519,
						Text:   "За всё, что Ты мне дал"},
					{
						Offset: 15.724137931,
						Length: 5.7931034483,
						Text:   "И за каждый новый вдох мой,"},
					{
						Offset: 21.5172413793,
						Length: 2.0606200734,
						Text:   "Буду петь"},
					{
						Offset: 23.5778614527,
						Length: 5.3876557887,
						Text:   "И на зака - те дня!"}}},
			{
				Name:     "Verse",
				Number:   2,
				Lang:     "R",
				Comments: []comments.Comment{},
				Text: []lyrics.TextLine{
					{
						Offset: -0.4763282443,
						Length: 2.7106906812,
						Text:   "Велик в любви"},
					{
						Offset: 2.2343624369,
						Length: 4.3863272183,
						Text:   "И даёшь мне милость"},
					{
						Offset: 6.6206896552,
						Length: 6.1750190239,
						Text:   "Ты справедлив и так благ ко мне"},
					{
						Offset: 12.7957086791,
						Length: 2.9284292519,
						Text:   "За все Твои дела"},
					{
						Offset: 15.724137931,
						Length: 5.7931034483,
						Text:   "Я хочу прославить"},
					{
						Offset: 21.5172413793,
						Length: 2.0606200734,
						Text:   "Чтоб петь Тебе"},
					{
						Offset: 23.5778614527,
						Length: 5.3876557887,
						Text:   "Есть десять тысяч причин!"}}},
			{
				Name:     "Verse",
				Number:   3,
				Lang:     "R",
				Comments: []comments.Comment{},
				Text: []lyrics.TextLine{
					{
						Offset: -0.4763282443,
						Length: 2.7106906812,
						Text:   "И когда"},
					{
						Offset: 2.2343624369,
						Length: 4.3863272183,
						Text:   "мои силы иссякнут"},
					{
						Offset: 6.6206896552,
						Length: 6.1750190239,
						Text:   "на закате дней, когда время придёт"},
					{
						Offset: 12.7957086791,
						Length: 2.9284292519,
						Text:   "Душа моя"},
					{
						Offset: 15.724137931,
						Length: 5.7931034483,
						Text:   "Тебя будет славить"},
					{
						Offset: 21.5172413793,
						Length: 2.0606200734,
						Text:   "Больше, чем даже"},
					{
						Offset: 23.5778614527,
						Length: 5.3876557887,
						Text:   "Десять ты - сяч лет!"}}},
			{
				Name:     "Chorus",
				Number:   0,
				Lang:     "R",
				Comments: []comments.Comment{},
				Text: []lyrics.TextLine{
					{
						Offset: -0.934376352,
						Length: 4.2447211796,
						Text:   "Славь душа Господа"},
					{
						Offset: 3.3103448276,
						Length: 3.3103448276,
						Text:   "Cлавь душа"},
					{
						Offset: 6.6206896552,
						Length: 5.6536462652,
						Text:   "Имя святое его"},
					{
						Offset: 12.2743359204,
						Length: 4.2773882175,
						Text:   "Пой Ему, как никогда"},
					{
						Offset: 16.5517241379,
						Length: 3.168300811,
						Text:   "Пой душа"},
					{
						Offset: 19.7200249489,
						Length: 9.2454922925,
						Text:   "Благослови Господа"}}}}}
	WhoYouSayIAm := lyrics.Lyrics{
		Name: "Who You Say I Am",
		Sections: []lyrics.Section{
			{
				Name:     "Verse",
				Number:   1,
				Lang:     "",
				Comments: []comments.Comment{},
				Text: []lyrics.TextLine{
					{
						Offset: 0.0157481992,
						Length: 3.8214611031,
						Text:   "Who am I that the highest King"},
					{
						Offset: 3.8372093023,
						Length: 4.8837209302,
						Text:   "Would welcome me"},
					{
						Offset: 8.7209302325,
						Length: 3.1395348838,
						Text:   "I was lost but He brought me in"},
					{
						Offset: 11.8604651163,
						Length: 3.8372093023,
						Text:   "Oh His love for me"}}},
			{
				Name:     "Chorus",
				Number:   1,
				Lang:     "",
				Comments: []comments.Comment{},
				Text: []lyrics.TextLine{
					{
						Offset: -0.6976744186,
						Length: 4.1860465116,
						Text:   "Who the Son sets free"},
					{
						Offset: 3.488372093,
						Length: 4.1860465117,
						Text:   "Oh is free indeed"},
					{
						Offset: 7.6744186047,
						Length: 8.0232558139,
						Text:   "I'm a child of God yes I am"}}},
			{
				Name:     "Verse",
				Number:   2,
				Lang:     "",
				Comments: []comments.Comment{},
				Text: []lyrics.TextLine{
					{
						Offset: 0.3488372093,
						Length: 3.4883720931,
						Text:   "Free at last He has ransomed me"},
					{
						Offset: 3.8372093024,
						Length: 4.8837209302,
						Text:   "His grace runs deep"},
					{
						Offset: 8.7209302326,
						Length: 3.1395348837,
						Text:   "While I was a slave to sin"},
					{
						Offset: 11.8604651163,
						Length: 3.8372093023,
						Text:   "Jesus died for me"}}},
			{
				Name:     "Chorus",
				Number:   2,
				Lang:     "",
				Comments: []comments.Comment{},
				Text: []lyrics.TextLine{
					{
						Offset: -0.6976744186,
						Length: 4.1860465116,
						Text:   "In my Father's house"},
					{
						Offset: 3.488372093,
						Length: 4.1860465116,
						Text:   "There's a place for me"},
					{
						Offset: 7.6744186046,
						Length: 8.023255814,
						Text:   "I'm a child of God yes I am"}}},
			{
				Name:     "Bridge",
				Number:   0,
				Lang:     "",
				Comments: []comments.Comment{},
				Text: []lyrics.TextLine{
					{
						Offset: -0.6976744186,
						Length: 4.5348837209,
						Text:   "I am chosen not forsaken"},
					{
						Offset: 3.8372093023,
						Length: 3.8372093023,
						Text:   "I am who You say I am"},
					{
						Offset: 7.6744186046,
						Length: 4.534883721,
						Text:   "You are for me not against me"},
					{
						Offset: 12.2093023256,
						Length: 3.488372093,
						Text:   "I am who You say I am"}}},
			{
				Name:     "Verse",
				Number:   1,
				Lang:     "R",
				Comments: []comments.Comment{},
				Text: []lyrics.TextLine{
					{
						Offset: 0.0157481992,
						Length: 3.8214611031,
						Text:   "Кто же я, что Всевышний Царь"},
					{
						Offset: 3.8372093023,
						Length: 4.8837209302,
						Text:   "Меня призвал"},
					{
						Offset: 8.7209302325,
						Length: 3.1395348838,
						Text:   "Я блуждал, но любовь Твоя"},
					{
						Offset: 11.8604651163,
						Length: 3.8372093023,
						Text:   "Меня нашла"}}},
			{
				Name:     "Verse",
				Number:   2,
				Lang:     "R",
				Comments: []comments.Comment{},
				Text: []lyrics.TextLine{
					{
						Offset: 0.3488372093,
						Length: 3.4883720931,
						Text:   "Я свободен искуплен Им"},
					{
						Offset: 3.8372093024,
						Length: 4.8837209302,
						Text:   "Он ми — ло — стив"},
					{
						Offset: 8.7209302326,
						Length: 3.1395348837,
						Text:   "Был я в рабстве греха, но Ты"},
					{
						Offset: 11.8604651163,
						Length: 3.8372093023,
						Text:   "Жизнью за — пла - тил"}}},
			{
				Name:     "Chorus",
				Number:   1,
				Lang:     "R",
				Comments: []comments.Comment{},
				Text: []lyrics.TextLine{
					{
						Offset: -0.6976744186,
						Length: 4.1860465116,
						Text:   "Дал свободу мне"},
					{
						Offset: 3.488372093,
						Length: 4.1860465117,
						Text:   "Иисус навек"},
					{
						Offset: 7.6744186047,
						Length: 8.0232558139,
						Text:   "Я дитя Ца - ря знаю я"}}},
			{
				Name:     "Chorus",
				Number:   2,
				Lang:     "R",
				Comments: []comments.Comment{},
				Text: []lyrics.TextLine{
					{
						Offset: -0.6976744186,
						Length: 4.1860465116,
						Text:   "В Доме есть всегда"},
					{
						Offset: 3.488372093,
						Length: 4.1860465116,
						Text:   "Место для меня"},
					{
						Offset: 7.6744186046,
						Length: 8.023255814,
						Text:   "Я дитя Ца - ря знаю я"}}},
			{
				Name:     "Bridge",
				Number:   0,
				Lang:     "R",
				Comments: []comments.Comment{},
				Text: []lyrics.TextLine{
					{
						Offset: -0.6976744186,
						Length: 4.5348837209,
						Text:   "Богом избран не оставлен"},
					{
						Offset: 3.8372093023,
						Length: 3.8372093023,
						Text:   "Я знаю, кто я в Тебе"},
					{
						Offset: 7.6744186046,
						Length: 4.534883721,
						Text:   "Твоим другом был я назван"},
					{
						Offset: 12.2093023256,
						Length: 3.488372093,
						Text:   "Я знаю, кто я в Тебе"}}}}}

	Tremble := lyrics.Lyrics{
		Name: "Tremble",
		Sections: []lyrics.Section{
			{
				Name:     "Verse",
				Number:   1,
				Lang:     "",
				Comments: []comments.Comment{},
				Text: []lyrics.TextLine{
					{
						Offset: 0,
						Length: 4.4594594594,
						Text:   "Peace, bring it all to peace"},
					{
						Offset: 4.4594594594,
						Length: 2.4324324325,
						Text:   "The storms surrounding me"},
					{
						Offset: 6.8918918919,
						Length: 5.6756756756,
						Text:   "Let it break at Your name"},
					{
						Offset: 12.5675675675,
						Length: 4.8648648649,
						Text:   "Still call the sea to still"},
					{
						Offset: 17.4324324324,
						Length: 2.4324324324,
						Text:   "The rage in me to still"},
					{
						Offset: 19.8648648648,
						Length: 5.6756756757,
						Text:   "Every wave at Your name"},
				},
			},
			{
				Name:     "Chorus",
				Number:   0,
				Lang:     "",
				Comments: []comments.Comment{},
				Text: []lyrics.TextLine{
					{
						Offset: -0.125749752,
						Length: 3.3689929952,
						Text:   "Jesus Jesus"},
					{
						Offset: 3.2432432432,
						Length: 3.0910167098,
						Text:   "You make the darkness tremble"},
					{
						Offset: 6.334259953,
						Length: 3.3954697767,
						Text:   "Jesus Jesus You silence fear"},
					{
						Offset: 9.7297297297,
						Length: 2.8437148543,
						Text:   "Jesus Jesus"},
					{
						Offset: 12.573444584,
						Length: 3.6427716322,
						Text:   "You make the darkness tremble"},
					{
						Offset: 16.2162162162,
						Length: 2.8378378379,
						Text:   "Jesus Jesus"}}},
			{
				Name:     "Verse",
				Number:   2,
				Lang:     "",
				Comments: []comments.Comment{},
				Text: []lyrics.TextLine{
					{
						Offset: 0,
						Length: 4.4594594594,
						Text:   "Breathe call these bones to live"},
					{
						Offset: 4.4594594594,
						Length: 2.4324324324,
						Text:   "Call these lungs to sing"},
					{
						Offset: 6.8918918918,
						Length: 5.6756756757,
						Text:   "Once again I will praise"}}},
			{
				Name:     "Bridge",
				Number:   1,
				Lang:     "",
				Comments: []comments.Comment{},
				Text: []lyrics.TextLine{
					{
						Offset: -0.8108108108,
						Length: 3.2432432432,
						Text:   "Your name is a light"},
					{
						Offset: 2.4324324324,
						Length: 3.2432432433,
						Text:   "That the shadows can't deny"},
					{
						Offset: 5.6756756757,
						Length: 6.0810810811,
						Text:   "Your name cannot be overcome"}}},
			{
				Name:     "Bridge",
				Number:   2,
				Lang:     "",
				Comments: []comments.Comment{},
				Text: []lyrics.TextLine{
					{
						Offset: -0.8108108108,
						Length: 3.2432432432,
						Text:   "Your name is alive forever lifted high"},
					{
						Offset: 2.4324324324,
						Length: 3.2432432433,
						Text:   "Your name cannot be overcome"}}},
			{
				Name:     "Tag",
				Number:   0,
				Lang:     "",
				Comments: []comments.Comment{},
				Text: []lyrics.TextLine{
					{
						Offset: 0,
						Length: 0,
						Text:   "Jesus Jesus"}}},
			{
				Name:     "Verse",
				Number:   1,
				Lang:     "R",
				Comments: []comments.Comment{},
				Text: []lyrics.TextLine{
					{
						Offset: 0,
						Length: 0,
						Text:   "Мир, принеси твой мир"},
					{
						Offset: 0,
						Length: 0,
						Text:   "Бури успокой"},
					{
						Offset: 0,
						Length: 0,
						Text:   "Повели им уйти!"},
					{
						Offset: 0,
						Length: 0,
						Text:   "Бог, слушают тебя"},
					{
						Offset: 0,
						Length: 0,
						Text:   "И бездны и моря"},
					{
						Offset: 0,
						Length: 0,
						Text:   "Голос твой так велик!"}}},
			{
				Name:     "Chorus",
				Number:   0,
				Lang:     "R",
				Comments: []comments.Comment{},
				Text: []lyrics.TextLine{
					{
						Offset: 0,
						Length: 0,
						Text:   "Иисус, иисус,"},
					{
						Offset: 0,
						Length: 0,
						Text:   "Приводишь тьму ты в трепет!"},
					{
						Offset: 0,
						Length: 0,
						Text:   "Иисус, иисус, и страх молчит!"},
					{
						Offset: 0,
						Length: 0,
						Text:   "Иисус, иисус,"},
					{
						Offset: 0,
						Length: 0,
						Text:   "Приводишь тьму ты в трепет!"},
					{
						Offset: 0,
						Length: 0,
						Text:   "Иисус, иисус"}}},
			{
				Name:     "Verse",
				Number:   2,
				Lang:     "R",
				Comments: []comments.Comment{},
				Text: []lyrics.TextLine{
					{
						Offset: 0,
						Length: 4.4594594594,
						Text:   "Жизнь даёшь сухим костям,"},
					{
						Offset: 4.4594594594,
						Length: 2.4324324324,
						Text:   "Даёшь дыханье мне"},
					{
						Offset: 6.8918918918,
						Length: 5.6756756757,
						Text:   "Что бы петь для тебя."}}},
			{
				Name:     "Bridge",
				Number:   1,
				Lang:     "R",
				Comments: []comments.Comment{},
				Text: []lyrics.TextLine{
					{
						Offset: -0.8108108108,
						Length: 3.2432432432,
						Text:   "Вся власть, сила вся"},
					{
						Offset: 2.4324324324,
						Length: 3.2432432433,
						Text:   "Только в имени твоём,"},
					{
						Offset: 5.6756756757,
						Length: 6.0810810811,
						Text:   "Весь мрак трепещет от тебя."}}},
			{
				Name:     "Bridge",
				Number:   2,
				Lang:     "R",
				Comments: []comments.Comment{},
				Text: []lyrics.TextLine{
					{
						Offset: -0.8108108108,
						Length: 3.2432432432,
						Text:   "Иисус, ты есть свет, Твоё имя на века"},
					{
						Offset: 2.4324324324,
						Length: 3.2432432433,
						Text:   "Никто не победит тебя."}}},
			{
				Name:     "Tag",
				Number:   0,
				Lang:     "R",
				Comments: []comments.Comment{},
				Text: []lyrics.TextLine{
					{
						Offset: 0,
						Length: 0,
						Text:   "Иисус, иисус,"}}}}}

	// fmt.Print(OneThingRemains.Name, WhoYouSayIAm.Name, BlessTheLord.Name)

	section := Tremble.Sections[6]
	// tests := OffsetMap{
	// 	{Offset: 0, Length: 4.4594594594},
	// 	{Offset: 4.4594594594, Length: 2.4324324325},
	// 	{Offset: 6.8918918919, Length: 5.6756756756},
	// 	{Offset: 12.5675675675, Length: 4.8648648649},
	// 	{Offset: 17.4324324324, Length: 2.4324324324},
	// 	{Offset: 19.8648648648, Length: 5.6756756757},
	// }
	//
	if section.Name != "Verse" || section.Number != 1 || section.Lang != "R" {

		t.Fatalf("expected section Verse 1 R, got %v", section.GetFullName())

	}

	// offsets := CalculateOffsets(Tremble, section)

	// for i, offset := range offsets {
	// 	if offset.Offset != tests[i].Offset {
	// 		t.Fatalf("test[%d] failed expected Offset=%v got=%v", i, tests[i].Offset, offset.Offset)
	// 	}
	//
	// 	if offset.Length != tests[i].Length {
	// 		t.Fatalf("test[%d] failed expected Length=%v got=%v", i, tests[i].Length, offset.Length)
	// 	}
	//
	// }

}
