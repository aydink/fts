package main

import (
	"fmt"
	"testing"
)

func TestPage(t *testing.T) {

	fmt.Println("------PAGE TEST-----------------")

	books := []Book{
		{Id: 1, Title: "Aşk tesadüfleri sever", Author: "Aşk Meşk", Type: "Roman", Genre: "Romantik", Year: 1900, NumPages: 120, Category: []string{"Yeni", "İndirim"}},
		{Id: 2, Title: "Aşk bunun neresinde, ey aşk", Author: "Ali Veli", Type: "Hikaye", Genre: "Deneme", Year: 2001, NumPages: 20, Category: []string{"Yeni", "Bedava"}},
		{Id: 3, Title: "Milletlerarası Özel Hukukta Kişilik Haklarının İnternet Yoluyla İhlalinde Sorumluluk", Author: "Esra Tekin", Type: "Kitap", Genre: "Akademik", Year: 2021, NumPages: 258, Category: []string{"Hukuk", "İnternet", "Özel Hukuk"}},
	}

	pages := []Page{
		{1, 1, "Devenin belini kiran son saman çöpüdür.", 1},
		{2, 1, "Bütün insanlar hür, haysiyet ve haklar bakimindan esit dogarlar. Akil ve vicdana sahiptirler ve birbirlerine karsi kardeslik zihniyeti ile hareket etmelidirler.", 2},
		{3, 1, "Bilmiyorum.", 3},
		{4, 1, "Mum kendiliginden söndü.", 4},
		{5, 1, "Mum kendi kendine söndü.", 5},
		{6, 1, "Babam bana yatakta kitap okumamami söyledi.", 6},
		{7, 1, "On, on bir, on iki, on üç, on dört, on bes, on alti, on yedi, on sekiz, on dokuz, yirmi.", 7},
		{8, 1, "Ek olarak yaslilar birbirleriyle sosyallesebilsin ve Amerikan hayatinin aktif üyeleri olarak kalabilsinler diye birçok topluluk kurulmustur.", 8},
		{9, 1, "Bazilari yalnizca zaman geçsin diye kitap okurlar.", 9},
		{10, 1, "Koyu kahverengi saçlari vardi.", 10},
		{11, 1, "Bu gemi okyanus yolculugu için uygun degil.", 11},
		{12, 1, "Bu kitap.", 12},
		{13, 1, "Hemen yolculuga hazirlan.", 13},
		{14, 1, "Masanin üstündeki hesap makinesi benim.", 14},
		{15, 1, "Simsek çakti.", 15},
		{16, 1, "Ask onu rüyalarinda görmektir.", 16},
		{17, 1, "Kimse benim fikirlerimi dinlemek istemiyor.", 17},
		{18, 1, "Sana satranç oynamayi ögretecegim.", 18},
		{19, 1, "Satranç oynamayi biliyor musun?", 19},
		{20, 1, "Bunlar çok eski kitaplar.", 20},
		{21, 1, "Bunlar benim kitaplarim.", 21},
		{22, 1, "Bunlar bizim kitaplarimiz.", 22},
		{23, 1, "Bunlar benim kalemlerim.", 23},
		{24, 1, "Bunlar her yerde satiliyor.", 24},
		{25, 1, "Bu kitap okul kütüphanesinin.", 25},
		{26, 1, "Köpekleri severim.", 26},
		{27, 1, "Bulasik makinesinin nasil çalistigini anlatabilir misin?", 27},
		{28, 1, "Zor durumlarla basa çikamiyor.", 28},
		{29, 1, "Günde en az yedi saat uyumak zorundayiz.", 29},
		{30, 1, "Sporu rekabet için degil zevk için yapiyorum.", 30},
		{31, 1, "Japonca konusamiyorum.", 31},
		{32, 1, "Yalnizca birkaç kisi vaktinde geldi.", 32},
		{33, 1, "Sadece birkaç kisi beni anladi.", 33},
		{34, 1, "Kosucuyum.", 34},
		{35, 1, "Geçen sene kurulan lunapark sagolsun sehir popüler oldu.", 35},
		{36, 1, "Onunla beraber oldugun sürece mutlu olamazsin.", 36},
		{37, 1, "Çocuklar yerde uyumak zorunda kalacaklar gibi.", 37},
		{38, 1, "Sarkilari gençler arasinda iyi biliniyor.", 38},
		{39, 1, "Bir davetiye aldim.", 39},
		{40, 1, "Matsuyama'da dogup büyüdüm.", 40},
		{41, 1, "Düsmanla anlasmaya vardilar.", 41},
		{42, 1, "Go büyük ihtimalle benim ülkemdeki en popüler Japon oyunu olsa da o bile bazi üniversite ögrencileri disinda pek bilinmiyor.", 42},
		{43, 1, "Mahjong taslari çok güzeller.", 43},
		{44, 1, "Mahjong genellikle dört kisi oynanan bir oyun.", 44},
		{45, 1, "Mahjong dünyada çok popüler olan oyunlardan biri.", 45},
		{46, 1, "Mahjong oynamayi biliyor musun?", 46},
		{47, 1, "Mahjong'u çok seviyorum.", 47},
		{48, 1, "Mahjong'da çok iyiymis.", 48},
		{49, 1, "Bu Mahjong.", 49},
		{50, 1, "Hayalim çok güçlü bir Mahjong oyuncusu olmak.", 50},
		{51, 1, "Japoncayi Japonya'da mahjong oynamak için ögreniyorum.", 51},
		{52, 1, "Atesin var mi?", 52},
		{53, 1, "Su köpek elimi isirmaya çalisti.", 53},
		{54, 1, "Evimin arkasinda bir kilise var.", 54},
		{55, 1, "Sam Tom'dan iki yas küçük.", 55},
		{56, 1, "Lütfen burayi imzalayin.", 56},
		{57, 1, "Zamanda geçmise seyahat etmenin imkansiz oldugu düsünülüyor.", 57},
		{58, 1, "Bu otobüs elli kisilik.", 58},
		{59, 1, "John cebinden bir anahtar çikardi.", 59},
		{60, 1, "John Floridali, karisi ise Kaliforniyali.", 60},
		{61, 1, "John, Fransizcayi iyi konusamiyor.", 61},
		{62, 1, "John birçok sise sarap içti.", 62},
		{63, 1, "John Bill'den daha zeki.", 63},
		{64, 1, "John Bill'in zayifligindan istifade etti.", 64},
		{65, 1, "John Bill kadar yasli degil; çok daha genç.", 65},
		{66, 1, "John o kadar telasliydi ki konusmaya vakti yoktu.", 66},
		{67, 1, "John ise asina.", 67},
		{68, 1, "John, New York'ta yasiyor.", 68},
		{69, 1, "John o kadar sesli konustu ki ikinci kattan bile duyabildim.", 69},
		{70, 1, "Is ben gelmeden önce bitmisti.", 70},
		{71, 1, "Birini taniyorum da ötekini degil.", 71},
		{72, 1, "Ailesini çok endiselendirdi.", 72},
		{73, 1, "Sadece beyaz kagit yeterli.", 73},
		{74, 1, "Üsüyüp isiticiyi açtim.", 74},
		{75, 1, "Dogdugum yer olan Nagasaki, güzel bir liman kentidir.", 75},
		{76, 1, "Kizin artik bir çocuk degil.", 76},
		{77, 1, "Hava bulutlaniyor.", 77},
		{78, 1, "Pardon ama radyoyu kisabilir misin acaba?", 78},
		{79, 1, "Tenis oynamak sagliklidir.", 79},
		{80, 1, "Saglikli olan adam sagligin degerini bilmez.", 80},
		{81, 1, "Tehlikenin farkinda olmayabilir.", 81},
		{82, 1, "Hangi dügmeye basacagimi söyler misin lütfen?", 82},
		{83, 1, "Onunla görüsebildigim için mutluyum.", 83},
		{84, 1, "Yakinlardaki küçük bir kasabada yasiyordu.", 84},
		{85, 1, "Kulübe katilmaz misin?", 85},
		{86, 1, "Kulübe katilmak isteyenler lütfen buraya isimlerini yazsin.", 86},
		{87, 1, "Kendimi nedense geceleri daha iyi hissediyorum.", 87},
		{88, 1, "Ummak bir strateji degildir.", 88},
		{89, 1, "Amcamlarda yedik.", 89},
		{90, 1, "Onu on dolara aldim.", 90},
		{91, 1, "Merhaba, nasilsin?", 91},
		{92, 1, "Iyi aksamlar Sikago.", 92},
		{93, 1, "Su kiz kim acaba.", 93},
		{94, 1, "Amerika'da arabalar yolun sag tarafini kullanirlar.", 94},
		{95, 1, "Yasli adam Kedi mi? diye sordu.", 95},
		{96, 1, "Bazen bir kiz miyim diye merak ediyorum.", 96},
		{97, 1, "Gerçekleri abartmayalim.", 97},
		{98, 1, "Deneyelim!", 98},
		{99, 1, "Üzülmeyin, her sey düzelecek.", 99},
		{100, 1, "Bu sadece bir sakaydi.", 100},
		{101, 2, "Üstü kalsin.", 1},
		{102, 2, "Bu oyunlar yetiskin kategorisi altinda listelenmis.", 2},
		{103, 2, "Tokyo'da.", 3},
		{104, 2, "Bu köprü on tondan fazlasini tasiyamaz.", 4},
		{105, 2, "Mum isigini severim.", 5},
		{106, 2, "Artik seni sevmiyorum.", 6},
		{107, 2, "Okulu sevmiyorum.", 7},
		{108, 2, "Rap sever misin?", 8},
		{109, 2, "Bu filmi sevdim.", 9},
		{110, 2, "Matematigi severim.", 10},
		{111, 2, "Sigara içmeyi biraktim.", 11},
		{112, 2, "Babana en iyi dileklerimle.", 12},
		{113, 2, "Tek basima yürüdüm.", 13},
		{114, 2, "Kapiyi kapatin, lütfen.", 14},
		{115, 2, "Erken Ortaçag cam üretimi Roma cam teknolojisinin bir devami midir?", 15},
		{116, 2, "Seni anlamak gerçekten çok zor.", 16},
		{117, 2, "Dün on yedinci dogum günümdü.", 17},
		{118, 2, "Evren nasil olustu?", 18},
		{119, 2, "Elektronik sigaradan çikan duman miktari hiç fena degil.", 19},
		{120, 2, "Kendimi kendi tanrim olarak görüyorum.", 20},
		{121, 2, "Eve dönerken yagmura yakalanmistim.", 21},
		{122, 2, "Tatoeba'ya yüzlerce cümle yazmak isterdim ama yapmam gereken seyler var.", 22},
		{123, 2, "Geç kaldigim için üzgünüm.", 23},
		{124, 2, "Açikça konusmak gerekirse, o güvenilmez biri.", 24},
		{125, 2, "Irene Pepperberg, Alex adini taktigi bir papaganin önünde yuvarlak bir tepsi tutuyor.", 25},
		{126, 2, "Japonya'da dogmus olmayi tercih ederdim.", 26},
		{127, 2, "Kizi onunla her yere gitmeye hevesli.", 27},
		{128, 2, "Biraz sakinlesmelisin.", 28},
		{129, 2, "Sen olmasaydin, o hâlâ hayatta olacakti.", 29},
		{130, 2, "Bana gülümsedi.", 30},
		{131, 2, "Çok daha iyi hissediyorum.", 31},
		{132, 2, "Türkiye'den döndükten sonra Türkçem zayifladi.", 32},
		{133, 2, "Politik dünyada pek çok düsmani var.", 33},
		{134, 2, "Onunla kahve dükkaninda bulusmaya söz verdi.", 34},
		{135, 2, "Büyük bedenimiz var, ama o renk mevcut degil.", 35},
		{136, 2, "Onu Jim diye çagirirlar.", 36},
		{137, 2, "Bence yarin yagmur yagmayacak.", 37},
		{138, 2, "Bes köfte istiyorum.", 38},
		{139, 2, "O kadar kötü birisi ki kimse ondan hoslanmaz.", 39},
		{140, 2, "Bir gün için bu kadari yeterli.", 40},
		{141, 2, "Artik çocuk degilim.", 41},
		{142, 2, "Daha dikkatli sür, aksi halde basin belaya girecek.", 42},
		{143, 2, "Alkolsüz içecekleriniz var mi ?", 43},
		{144, 2, "Çocuklarin daha çok uykuya ihtiyaci vardir.", 44},
		{145, 2, "Odasina girdi.", 45},
		{146, 2, "Kameram Nikon'dur.", 46},
		{147, 2, "Hiç süphe yok ki Ingilizce dünyada en çok konusulan dildir.", 47},
		{148, 2, "Doktor olacak.", 48},
		{149, 2, "Yumi ögretmen olacak.", 49},
		{150, 2, "Bu tirtil harika bir kelebek olacak.", 50},
		{151, 2, "Bu tirtil harika bir kelebege dönüsecek.", 51},
		{152, 2, "Bugün hava kötü.", 52},
		{153, 2, "Beklemekten baska çare yoktu.", 53},
		{154, 2, "Odada 2 tane pencere var.", 54},
		{155, 2, "Onu tanidikça daha çok seversin.", 55},
		{156, 2, "Hastaymis gibi gözüküyor.", 56},
		{157, 2, "Çok lezzetli.", 57},
		{158, 2, "Bu kiyafetim çok demode.", 58},
		{159, 2, "Sonsuza dek burada kalamam.", 59},
		{160, 2, "Hos geldin!", 60},
		{161, 2, "Nasilsin?", 61},
		{162, 2, "Sevebilirim.", 62},
		{163, 2, "Rengin solmus.", 63},
		{164, 2, "Seni nasil da özledim!", 64},
		{165, 2, "Az daha treni kaçiriyordum.", 65},
		{166, 2, "Bir gece daha kalmak istiyorum. Mümkün mü?", 66},
		{167, 2, "Mahjong en ilginç oyunlardan biri.", 67},
		{168, 2, "Japonya'ya gidip Mahjong'da Japonlari yenmek istiyorum.", 68},
		{169, 2, "Sirada bir elma var.", 69},
		{170, 2, "Atlas Okyanusu Amerika'yi Avrupa'dan ayirir.", 70},
		{171, 2, "Sikildim.", 71},
		{172, 2, "Evliyim ve iki çocugum var.", 72},
		{173, 2, "Haberler onu üzdü.", 73},
		{174, 2, "Çok sarap içmiyorum.", 74},
		{175, 2, "Nancy onu bir partiye çagirdi.", 75},
		{176, 2, "Çok fazla yersen sismanlarsin.", 76},
		{177, 2, "Golfün büyük bir hayraniyim.", 77},
		{178, 2, "Kötü mü?", 78},
		{179, 2, "Insanlar savastan korkar.", 79},
		{180, 2, "Iyi hissetmiyorum.", 80},
		{181, 2, "Ögretmenimiz bize çok ödev verir.", 81},
		{182, 2, "Tesekkür ederim!", 82},
		{183, 2, "Tesekkür ederiz!", 83},
		{184, 2, "Çok tesekkür ederim!", 84},
		{185, 2, "Çok tesekkür ederiz!", 85},
		{186, 2, "Çok tesekkürler!", 86},
		{187, 2, "Iyi bir doktor kendi yöntemlerini uygular.", 87},
		{188, 2, "O, patronunu dinlememe numarasi yapti.", 88},
		{189, 2, "Onu Kaliforniya'ya gönderiyorum.", 89},
		{190, 2, "Adam burada.", 90},
		{191, 2, "Adami görüyorum.", 91},
		{192, 2, "Kadinin fotografini çekiyorum.", 92},
		{193, 2, "Kadin okuyor.", 93},
		{194, 2, "Pekin, Çin'in baskentidir.", 94},
		{195, 2, "Bu elma tatli.", 95},
		{196, 2, "Kaç tane elma var?", 96},
		{197, 2, "Seni seviyorum.", 97},
		{198, 2, "Almanya parlamenter bir cumhuriyettir.", 98},
		{199, 2, "Internette Tatar dilinde çok az site vardir.", 99},
		{200, 2, "Hiç kimse bilmiyor.", 100},
		{201, 3, "Sigara içmek size çok zarar verecektir.", 1},
		{202, 3, "Eve hos geldin.", 2},
		{203, 3, "O okulunu samimiyetle seviyor.", 3},
		{204, 3, "Bir sey degil.", 4},
		{205, 3, "Eve gidebilir miyiz?", 5},
		{206, 3, "Doktor olarak iyi degil.", 6},
		{207, 3, "Ne is yapiyorsun?", 7},
		{208, 3, "Ingilizce konusabiliyor musun?", 8},
		{209, 3, "Kedi uzaktayken fareler oynayacak.", 9},
		{210, 3, "Demiryolu istasyonu nerede?", 10},
		{211, 3, "Tanistigimiza memnun oldum.", 11},
		{212, 3, "Yakinda o bir baba olacak.", 12},
		{213, 3, "Çocuklarin çogunlugu degisimi çok iyi duyamazlar.", 13},
		{214, 3, "Kredi kartiyla ödeyebilir miyim?", 14},
		{215, 3, "Bana yardim edebilir misin?", 15},
		{216, 3, "Bir, iki, üç, dört, bes, alti, yedi, sekiz, dokuz, on.", 16},
		{217, 3, "Susadim.", 17},
		{218, 3, "Iyi aksamlar. Nasilsin?", 18},
		{219, 3, "Iyi aksamlar. Nasilsiniz?", 19},
		{220, 3, "Her parlayan sey altin degildir.", 20},
		{221, 3, "Bir kilicim yok.", 21},
		{222, 3, "Singapurluyum.", 22},
		{223, 3, "Banyo nerede?", 23},
		{224, 3, "Okulun nerede?", 24},
		{225, 3, "Iyi aksamlar.", 25},
		{226, 3, "Büyükbaban nerede yasiyor?", 26},
		{227, 3, "Deden nerede yasiyor?", 27},
		{228, 3, "Asagidaki sorulari Ingilizce olarak cevapla.", 28},
		{229, 3, "Son olarak o Amerika'ya gitti.", 29},
		{230, 3, "Balik tutmayi seviyorum.", 30},
		{231, 3, "Dün sicakti.", 31},
		{232, 3, "Dün orasi sicakti.", 32},
		{233, 3, "Mademki bos vaktimiz var, o zaman sinemaya gidelim.", 33},
		{234, 3, "Saatlerdir bekliyorum.", 34},
		{235, 3, "Adam gibi davran.", 35},
		{236, 3, "Faturayi ödemedigi için suyu kestiler.", 36},
		{237, 3, "Güller açiyor.", 37},
		{238, 3, "Her kimin ihtiyaci olursa ona yardim et.", 38},
		{239, 3, "Evim otobüs duragina yakin.", 39},
		{240, 3, "Zamanin ölçüsü nedir?", 40},
		{241, 3, "Buraya dün aksam altida geldik.", 41},
		{242, 3, "Ay bulutlarin üzerinde kaldi.", 42},
		{243, 3, "Seni çok seviyorum.", 43},
		{244, 3, "Kisa saç stilini severim.", 44},
		{245, 3, "Konuya Fransiz kaldim.", 45},
		{246, 3, "Babanin nereye gittigini biliyor musun?", 46},
		{247, 3, "Benim adim Edgar Degas.", 47},
		{248, 3, "Maria kiyafete çok para harciyor.", 48},
		{249, 3, "Havuçlari tencereye koy.", 49},
		{250, 3, "Yarin burada olacagim.", 50},
		{251, 3, "Haydi Japonya'yi yenelim!", 51},
		{252, 3, "Siz burada bir ögretmen mi, yoksa ögrenci misiniz?", 52},
		{253, 3, "Siz bir ögretmen misiniz? Evet, ögretmenim.", 53},
		{254, 3, "Sizin bir ögretmen oldugunuzu biliyorum.", 54},
		{255, 3, "Beni liderinize götürün.", 55},
		{256, 3, "Adin ne?", 56},
		{257, 3, "Ben Fransizim.", 57},
		{258, 3, "Ben saglikliyim.", 58},
		{259, 3, "Elmalari sever misin?", 59},
		{260, 3, "Kabul ediyorum.", 60},
		{261, 3, "Ögretiyorum.", 61},
		{262, 3, "Ne?", 62},
		{263, 3, "Benim adim Andrea.", 63},
		{264, 3, "Kaç yasindasin?", 64},
		{265, 3, "Benim basim agriyor.", 65},
		{266, 3, "O çok pahali!", 66},
		{267, 3, "Hava soguk.", 67},
		{268, 3, "Dogru söylüyorsun.", 68},
		{269, 3, "Haklisin.", 69},
		{270, 3, "Ben çok yorgunum.", 70},
		{271, 3, "Her gün Ingilizce çalisiyor musun?", 71},
		{272, 3, "Bizimle burada kal.", 72},
		{273, 3, "Lütfen cevaplayin.", 73},
		{274, 3, "Ben dürüst bir insanim.", 74},
		{275, 3, "Isigi kapatir misiniz?", 75},
		{276, 3, "Almanca konusuyor musun?", 76},
		{277, 3, "Iste size bir mektup.", 77},
		{278, 3, "Merhaba, Tom.", 78},
		{279, 3, "Yakinda görüsürüz!", 79},
		{280, 3, "Ne zaman size yazilmis ve anlamadiginiz bir seyiniz varsa, ne yapabileceksiniz, ya beklenmedik sonuçlar alirsaniz?", 80},
		{281, 3, "Bu kisi kim?", 81},
		{282, 3, "O kitabi hiç okumadim.", 82},
		{283, 3, "Seninle birlikte mi gitmeliyim?", 83},
		{284, 3, "Onun 100 dolardan az parasi yok.", 84},
		{285, 3, "O iyi bir yüzücüdür.", 85},
		{286, 3, "Fileyle kelebek yakaladim.", 86},
		{287, 3, "Benim çok iyi bir sözlügüm yok.", 87},
		{288, 3, "Istasyona giderken ben seni geçtim.", 88},
		{289, 3, "Dilinizi anlayabiliyorum.", 89},
		{290, 3, "Seninle seyahat etmek istiyorum.", 90},
		{291, 3, "Iki yildir ilk defa bir film izledim.", 91},
		{292, 3, "Ben bekarim.", 92},
		{293, 3, "Ben evdeyim.", 93},
		{294, 3, "Londra'dayken Mary ve John'la karsilastim.", 94},
		{295, 3, "Ben Ken'e inaniyorum.", 95},
		{296, 3, "Bu grupla kendimi tanitmak istemiyorum.", 96},
		{297, 3, "Ben senden daha güzelim.", 97},
		{298, 3, "Bana öyle görünüyor ki sen hatalisin.", 98},
		{299, 3, "Ve onu üç günde tekrar kaldiracagim.", 99},
		{300, 3, "Amcamin cadde boyunca bir magazasi var.", 100},
		{301, 3, "Auckland, Yeni Zelanda'da bir sehirdir.", 101},
		{302, 3, "Elmalar kirmizi veya yesildir.", 102},
		{303, 3, "Batman, Robin ile arkadastir.", 103},
		{304, 3, "Auckland'in bir milyon nüfusu vardir.", 104},
		{305, 3, "Elmalar kirmizidir.", 105},
		{306, 3, "Bu bir sürpriz.", 106},
		{307, 3, "Hava yagmurlu.", 107},
		{308, 3, "Bu önemli degil.", 108},
		{309, 3, "Bileti unutma.", 109},
		{310, 3, "Sizi tekrar görmekten memnunum.", 110},
		{311, 3, "Bir seyi degistirmeyecek.", 111},
		{312, 3, "O bir seyi degistirmeyecek.", 112},
		{313, 3, "Tanriya sükür.", 113},
		{314, 3, "Günaydin, Mike.", 114},
		{315, 3, "Iyi uyu, Timmy.", 115},
		{316, 3, "Mutlu yillar Muiriel!", 116},
		{317, 3, "Tesekkürler. Bir sey degil.", 117},
		{318, 3, "Hepinize günaydin.", 118},
		{319, 3, "O bana hirsizligin ne kadar yanlis bir sey oldugunu anlatti.", 119},
		{320, 3, "O benim daga tek basima tirmanmamin imkansiz oldugunu düsünüyor.", 120},
		{321, 3, "Bütün serveti ve söhretine ragmen, o mutsuz.", 121},
		{322, 3, "Ne yapacagimi bilmiyorum.", 122},
		{323, 3, "Betty klasik müzigi sever.", 123},
		{324, 3, "Kâgidin var mi?", 124},
		{325, 3, "Restoranimiz Güney Otogari'na yakin.", 125},
		{326, 3, "Yemegimiz ucuz.", 126},
		{327, 3, "Bugün tek basina mi geldin?", 127},
		{328, 3, "Yeni misin?", 128},
		{329, 3, "Etli pilav sekiz yuan. Vejetaryen pilav sadece dört yuan.", 129},
		{330, 3, "Fare burada! Git de kediyi çagir!", 130},
		{331, 3, "Asçi nerede?", 131},
		{332, 3, "Hollanda küçük bir ülkedir.", 132},
		{333, 3, "O bir kitap okuyor mu? Evet, o okuyor.", 133},
		{334, 3, "Esperantoyu yayin!", 134},
		{335, 3, "Burada biri var mi?", 135},
		{336, 3, "Kitaplari bu ögrenciye verdim.", 136},
		{337, 3, "Ne yapabilirim?", 137},
		{338, 3, "Merhaba! Nasilsin?", 138},
		{339, 3, "Lütfen bana nerede yasayacagini söyle.", 139},
		{340, 3, "Ne oldugunu biliyor musun?", 140},
		{341, 3, "Ingilizce benim ana dilim.", 141},
		{342, 3, "Bu renk hosuna gidiyor mu?", 142},
		{343, 3, "Gök mavidir.", 143},
		{344, 3, "Bulutlar mavi gökte yüzüyor.", 144},
		{345, 3, "Onun bisikleti mavi.", 145},
		{346, 3, "Benim gözlerim mavi.", 146},
		{347, 3, "Bugün ya da yarin gitmen pek fark yaratmayacak.", 147},
		{348, 3, "Mohan ile top oynamaya gidiyorum.", 148},
		{349, 3, "Everest Dagi dünyanin en yüksek zirvesidir.", 149},
		{350, 3, "Nerede oturmak istiyorsun?", 150},
		{351, 3, "Bir zaman makinen oldugunu hayal et.", 151},
		{352, 3, "Ne zaman geri döneceksin?", 152},
		{353, 3, "O, bankta arkadasiyla oturuyor.", 153},
		{354, 3, "O her sabah kosmaya gider.", 154},
		{355, 3, "Onun birçok tarih kitabi var.", 155},
		{356, 3, "Bir köpegim var.", 156},
		{357, 3, "Kentteki en iyi otel hangisi?", 157},
		{358, 3, "Benim adim Wang.", 158},
		{359, 3, "Twitter kullaniyorum.", 159},
		{360, 3, "Ted trompet çalmayi sever.", 160},
		{361, 3, "Buna ihtiyacim var.", 161},
		{362, 3, "Bozuk kamerayi buldum.", 162},
		{363, 3, "Kendimi tanitmama izin ver.", 163},
		{364, 3, "Benim adim Ludwig.", 164},
		{365, 3, "Ben Anton.", 165},
		{366, 3, "Zamanim yok.", 166},
		{367, 3, "Kimildama.", 167},
		{368, 3, "Sekreter mektubu bir zarfa yerlestirdi.", 168},
		{369, 3, "Onun iyi bir piyanist oldugunu söylemeye gerek yok", 169},
		{370, 3, "Tarih çalismayi severim.", 170},
		{371, 3, "Niye o uzaga kostu?", 171},
		{372, 3, "Mary zaten basladi.", 172},
		{373, 3, "Çogu erkek çocugu bilgisayar oyunlarini sever.", 173},
		{374, 3, "Dün cumartesi degil, pazardi.", 174},
		{375, 3, "Ben hamileyim.", 175},
		{376, 3, "Ben tenis kulübünün bir üyesiyim.", 176},
		{377, 3, "Örgütümüze nasil katkida bulunabilirsiniz?", 177},
		{378, 3, "1956'da Atina'da dogdu.", 178},
		{379, 3, "Hava bugün sicak.", 179},
		{380, 3, "Orasi Sirbistan'in üçüncü büyük sehridir.", 180},
		{381, 3, "Güzel soru.", 181},
		{382, 3, "Geç oldu.", 182},
		{383, 3, "O yepyeni.", 183},
		{384, 3, "Hesap lütfen.", 184},
		{385, 3, "Basim agriyor.", 185},
		{386, 3, "Hâlâ Çinceyi iyi yazamiyorum.", 186},
		{387, 3, "Çok isime yaradi.", 187},
		{388, 3, "Ne adi bir davranis!", 188},
		{389, 3, "Yurtdisinda okuma kararim ebeveynlerimi sasirtti.", 189},
		{390, 3, "O ne zaman dogdu?", 190},
		{391, 3, "Tüm bu yillari kaybettin.", 191},
		{392, 3, "Lütfen beni affet.", 192},
		{393, 3, "Kusura bakmayin, ben kayboldum.", 193},
		{394, 3, "Kayip misin?", 194},
		{395, 3, "Benim bir sorum var.", 195},
		{396, 3, "Tekrar deneyecegim.", 196},
		{397, 3, "Mutluluk nedir?", 197},
		{398, 3, "Nerede yasiyorsun?", 198},
		{399, 3, "Sen kimsin?", 199},
		{400, 3, "Kore'de hangi diller konusuluyor?", 200},
	}

	analyzer := NewSimpleAnalyzer(NewSimpleTokenizer())

	turkishFilter := NewTurkishLowercaseFilter()
	turkishAccentFilter := NewTurkishAccentFilter()

	analyzer.AddTokenFilter(turkishFilter)
	analyzer.AddTokenFilter(turkishAccentFilter)

	index := NewPageIndex(analyzer)

	for _, page := range pages {
		index.Add(&page)
	}

	//fmt.Println("Number of books:", len(books))
	//fmt.Printf("%+v\n", index)

	q := "çok"

	fmt.Printf("\nsearching for: %s\n", q)
	fmt.Println("----------------------------")
	fmt.Printf("%+v\n", index.Search(q))

	hits := index.Search(q)

	for _, hit := range hits {

		fmt.Println(index.bookId[hit.docId], books[pages[hit.docId].BookId-1].Title)
		fmt.Println("---", pages[hit.docId].Content)
	}

}
