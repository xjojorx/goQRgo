package main

import (
	"encoding/binary"
	"fmt"
	"math/bits"
	"slices"
	"unicode"
)

type EncodingMode int
const (
  Numeric EncodingMode = iota // 0001
  Alphanumeric // 0010
  Byte // 0100
  Kanji // 1000
)

type CorrectionLevel rune
const (
  CorrectionL CorrectionLevel = 'L'
  CorrectionM CorrectionLevel = 'M'
  CorrectionQ CorrectionLevel = 'Q'
  CorrectionH CorrectionLevel = 'H'
)

type Version struct {
  nversion int
  correction CorrectionLevel
  capNum int
  capAlpha int
  capByte int
  capKanji int
  totalWords int
}

func (v Version) CharCountLength(mode EncodingMode) int {
  if v.nversion < 10 {
    switch mode {
    case Numeric:
      return 10
    case Alphanumeric:
      return 9
    case Byte:
      return 8
    case Kanji:
      return 8
    }
  } else if v.nversion < 27 {
    switch mode {
    case Numeric:
      return 12
    case Alphanumeric:
      return 11
    case Byte:
      return 16
    case Kanji:
      return 10
    }
  } else {
    switch mode {
    case Numeric:
      return 14
    case Alphanumeric:
      return 13
    case Byte:
      return 16
    case Kanji:
      return 12
    }
  }
  return 0
}

func main() {
  input := "HELLO WORLD"
  corrLvl := CorrectionQ //should read from args

  mode := encodingFormat(input)
  fmt.Printf("mode: %#v\n", mode)

  version := determineVersion(input, corrLvl, mode)

  fmt.Printf("input: '%s', mode: '%d', correction: '%s', version: '%d'\n", input, mode, string(corrLvl), version.nversion)

  encoded := encode(input, version, mode, corrLvl)

  fmt.Printf("encoded as: %08b\n", encoded)

}

func encodingFormat(input string) EncodingMode {
  alpha_symbols := []rune{'$','%','*','+','-','.','/',':',' '}
  mode := Numeric

  for _, char := range input {
    if !unicode.IsDigit(char) {
      //is not numeric then
      if unicode.IsUpper(char) || slices.Contains(alpha_symbols, char) {
        mode = Alphanumeric
      } else {
        //kanji or byte
        if unicode.Is(unicode.Han, char) {
          mode = Kanji
        } else {
          mode = Byte
          break
        }
      }
    }
  }

  return mode
}

func listVersions() []Version {
  return []Version{
    {nversion: 1, correction: CorrectionL, capNum: 41, capAlpha: 25, capByte: 17, capKanji: 10, totalWords: 19},
    {nversion: 1, correction: CorrectionM, capNum: 34, capAlpha: 20, capByte: 14, capKanji: 8, totalWords: 16},
    {nversion: 1, correction: CorrectionQ, capNum: 27, capAlpha: 16, capByte: 11, capKanji: 7, totalWords: 13},
    {nversion: 1, correction: CorrectionH, capNum: 17, capAlpha: 10, capByte: 7, capKanji: 4, totalWords: 9},

    {nversion: 2 ,correction: CorrectionL, capNum: 77, capAlpha: 47, capByte: 32, capKanji: 20, totalWords: 34},
    {nversion: 2 ,correction: CorrectionM, capNum: 63, capAlpha: 38, capByte: 26, capKanji: 16, totalWords: 28},
    {nversion: 2 ,correction: CorrectionQ, capNum: 48, capAlpha: 29, capByte: 20, capKanji: 12, totalWords: 22},
    {nversion: 2 ,correction: CorrectionH, capNum: 34, capAlpha: 20, capByte: 14, capKanji: 8, totalWords: 16},

    {nversion: 3 ,correction: CorrectionL, capNum: 127, capAlpha: 77, capByte: 53, capKanji: 32, totalWords: 55},
    {nversion: 3 ,correction: CorrectionM, capNum: 101, capAlpha: 61, capByte: 42, capKanji: 26, totalWords: 44},
    {nversion: 3 ,correction: CorrectionQ, capNum: 77, capAlpha: 47, capByte: 32, capKanji: 20, totalWords: 34},
    {nversion: 3 ,correction: CorrectionH, capNum: 58, capAlpha: 35, capByte: 24, capKanji: 15, totalWords: 26},

    {nversion: 4 ,correction: CorrectionL, capNum: 187, capAlpha: 114, capByte: 78, capKanji: 48, totalWords: 80},
    {nversion: 4 ,correction: CorrectionM, capNum: 149, capAlpha: 90, capByte: 62, capKanji: 38, totalWords: 64},
    {nversion: 4 ,correction: CorrectionQ, capNum: 111, capAlpha: 67, capByte: 46, capKanji: 28, totalWords: 48},
    {nversion: 4 ,correction: CorrectionH, capNum: 82, capAlpha: 50, capByte: 34, capKanji: 21, totalWords: 36},

    {nversion: 5 ,correction: CorrectionL, capNum: 255, capAlpha: 154, capByte: 106, capKanji: 65, totalWords: 108},
    {nversion: 5 ,correction: CorrectionM, capNum: 202, capAlpha: 122, capByte: 84, capKanji: 52, totalWords: 86},
    {nversion: 5 ,correction: CorrectionQ, capNum: 144, capAlpha: 87, capByte: 60, capKanji: 37, totalWords: 62},
    {nversion: 5 ,correction: CorrectionH, capNum: 106, capAlpha: 64, capByte: 44, capKanji: 27, totalWords: 46},

    {nversion: 6 ,correction: CorrectionL, capNum: 322, capAlpha: 195, capByte: 134, capKanji: 82, totalWords: 136},
    {nversion: 6 ,correction: CorrectionM, capNum: 255, capAlpha: 154, capByte: 106, capKanji: 65, totalWords: 108},
    {nversion: 6 ,correction: CorrectionQ, capNum: 178, capAlpha: 108, capByte: 74, capKanji: 45, totalWords: 76},
    {nversion: 6 ,correction: CorrectionH, capNum: 139, capAlpha: 84, capByte: 58, capKanji: 36, totalWords: 60},

    {nversion: 7 ,correction: CorrectionL, capNum: 370, capAlpha: 224, capByte: 154, capKanji: 95, totalWords: 156},
    {nversion: 7 ,correction: CorrectionM, capNum: 293, capAlpha: 178, capByte: 122, capKanji: 75, totalWords: 124},
    {nversion: 7 ,correction: CorrectionQ, capNum: 207, capAlpha: 125, capByte: 86, capKanji: 53, totalWords: 88},
    {nversion: 7 ,correction: CorrectionH, capNum: 154, capAlpha: 93, capByte: 64, capKanji: 39, totalWords: 66},

    {nversion: 8 ,correction: CorrectionL, capNum: 461, capAlpha: 279, capByte: 192, capKanji: 118, totalWords: 194},
    {nversion: 8 ,correction: CorrectionM, capNum: 365, capAlpha: 221, capByte: 152, capKanji: 93, totalWords: 154},
    {nversion: 8 ,correction: CorrectionQ, capNum: 259, capAlpha: 157, capByte: 108, capKanji: 66, totalWords: 110},
    {nversion: 8 ,correction: CorrectionH, capNum: 202, capAlpha: 122, capByte: 84, capKanji: 52, totalWords: 86},

    {nversion: 9 ,correction: CorrectionL, capNum: 552, capAlpha: 335, capByte: 230, capKanji: 141, totalWords: 232},
    {nversion: 9 ,correction: CorrectionM, capNum: 432, capAlpha: 262, capByte: 180, capKanji: 111, totalWords: 182},
    {nversion: 9 ,correction: CorrectionQ, capNum: 312, capAlpha: 189, capByte: 130, capKanji: 80, totalWords: 132},
    {nversion: 9 ,correction: CorrectionH, capNum: 235, capAlpha: 143, capByte: 98, capKanji: 60, totalWords: 100},

    {nversion: 10 ,correction: CorrectionL, capNum: 652, capAlpha: 395, capByte: 271, capKanji: 167, totalWords: 274},
    {nversion: 10 ,correction: CorrectionM, capNum: 513, capAlpha: 311, capByte: 213, capKanji: 131, totalWords: 216},
    {nversion: 10 ,correction: CorrectionQ, capNum: 364, capAlpha: 221, capByte: 151, capKanji: 93, totalWords: 154},
    {nversion: 10 ,correction: CorrectionH, capNum: 288, capAlpha: 174, capByte: 119, capKanji: 74, totalWords: 122},

    {nversion: 11 ,correction: CorrectionL, capNum: 772, capAlpha: 468, capByte: 321, capKanji: 198, totalWords: 324},
    {nversion: 11 ,correction: CorrectionM, capNum: 604, capAlpha: 366, capByte: 251, capKanji: 155, totalWords: 254},
    {nversion: 11 ,correction: CorrectionQ, capNum: 427, capAlpha: 259, capByte: 177, capKanji: 109, totalWords: 180},
    {nversion: 11 ,correction: CorrectionH, capNum: 331, capAlpha: 200, capByte: 137, capKanji: 85, totalWords: 140},

    {nversion: 12 ,correction: CorrectionL, capNum: 883, capAlpha: 535, capByte: 367, capKanji: 226, totalWords: 370},
    {nversion: 12 ,correction: CorrectionM, capNum: 691, capAlpha: 419, capByte: 287, capKanji: 177, totalWords: 290},
    {nversion: 12 ,correction: CorrectionQ, capNum: 489, capAlpha: 296, capByte: 203, capKanji: 125, totalWords: 206},
    {nversion: 12 ,correction: CorrectionH, capNum: 374, capAlpha: 227, capByte: 155, capKanji: 96, totalWords: 158},

    {nversion: 13 ,correction: CorrectionL, capNum: 1022, capAlpha: 619, capByte: 425, capKanji: 262, totalWords: 428},
    {nversion: 13 ,correction: CorrectionM, capNum: 796, capAlpha: 483, capByte: 331, capKanji: 204, totalWords: 334},
    {nversion: 13 ,correction: CorrectionQ, capNum: 580, capAlpha: 352, capByte: 241, capKanji: 149, totalWords: 244},
    {nversion: 13 ,correction: CorrectionH, capNum: 427, capAlpha: 259, capByte: 177, capKanji: 109, totalWords: 180},

    {nversion: 14 ,correction: CorrectionL, capNum: 1101, capAlpha: 667, capByte: 458, capKanji: 282, totalWords: 461},
    {nversion: 14 ,correction: CorrectionM, capNum: 871, capAlpha: 528, capByte: 362, capKanji: 223, totalWords: 365},
    {nversion: 14 ,correction: CorrectionQ, capNum: 621, capAlpha: 376, capByte: 258, capKanji: 159, totalWords: 261},
    {nversion: 14 ,correction: CorrectionH, capNum: 468, capAlpha: 283, capByte: 194, capKanji: 120, totalWords: 197},

    {nversion: 15 ,correction: CorrectionL, capNum: 1250, capAlpha: 758, capByte: 520, capKanji: 320, totalWords: 523},
    {nversion: 15 ,correction: CorrectionM, capNum: 991, capAlpha: 600, capByte: 412, capKanji: 254, totalWords: 415},
    {nversion: 15 ,correction: CorrectionQ, capNum: 703, capAlpha: 426, capByte: 292, capKanji: 180, totalWords: 295},
    {nversion: 15 ,correction: CorrectionH, capNum: 530, capAlpha: 321, capByte: 220, capKanji: 136, totalWords: 223},

    {nversion: 16 ,correction: CorrectionL, capNum: 1408, capAlpha: 854, capByte: 586, capKanji: 361, totalWords: 589},
    {nversion: 16 ,correction: CorrectionM, capNum: 1082, capAlpha: 656, capByte: 450, capKanji: 277, totalWords: 453},
    {nversion: 16 ,correction: CorrectionQ, capNum: 775, capAlpha: 470, capByte: 322, capKanji: 198, totalWords: 325},
    {nversion: 16 ,correction: CorrectionH, capNum: 602, capAlpha: 365, capByte: 250, capKanji: 154, totalWords: 253},

    {nversion: 17 ,correction: CorrectionL, capNum: 1548, capAlpha: 938, capByte: 644, capKanji: 397, totalWords: 647},
    {nversion: 17 ,correction: CorrectionM, capNum: 1212, capAlpha: 734, capByte: 504, capKanji: 310, totalWords: 507},
    {nversion: 17 ,correction: CorrectionQ, capNum: 876, capAlpha: 531, capByte: 364, capKanji: 224, totalWords: 367},
    {nversion: 17 ,correction: CorrectionH, capNum: 674, capAlpha: 408, capByte: 280, capKanji: 173, totalWords: 283},

    {nversion: 18 ,correction: CorrectionL, capNum: 1725, capAlpha: 1046, capByte: 718, capKanji: 442, totalWords: 721},
    {nversion: 18 ,correction: CorrectionM, capNum: 1346, capAlpha: 816, capByte: 560, capKanji: 345, totalWords: 563},
    {nversion: 18 ,correction: CorrectionQ, capNum: 948, capAlpha: 574, capByte: 394, capKanji: 243, totalWords: 397},
    {nversion: 18 ,correction: CorrectionH, capNum: 746, capAlpha: 452, capByte: 310, capKanji: 191, totalWords: 313},

    {nversion: 19 ,correction: CorrectionL, capNum: 1903, capAlpha: 1153, capByte: 792, capKanji: 488, totalWords: 795},
    {nversion: 19 ,correction: CorrectionM, capNum: 1500, capAlpha: 909, capByte: 624, capKanji: 384, totalWords: 627},
    {nversion: 19 ,correction: CorrectionQ, capNum: 1063, capAlpha: 644, capByte: 442, capKanji: 272, totalWords: 445},
    {nversion: 19 ,correction: CorrectionH, capNum: 813, capAlpha: 493, capByte: 338, capKanji: 208, totalWords: 341},

    {nversion: 20 ,correction: CorrectionL, capNum: 2061, capAlpha: 1249, capByte: 858, capKanji: 528, totalWords: 861},
    {nversion: 20 ,correction: CorrectionM, capNum: 1600, capAlpha: 970, capByte: 666, capKanji: 410, totalWords: 669},
    {nversion: 20 ,correction: CorrectionQ, capNum: 1159, capAlpha: 702, capByte: 482, capKanji: 297, totalWords: 485},
    {nversion: 20 ,correction: CorrectionH, capNum: 919, capAlpha: 557, capByte: 382, capKanji: 235, totalWords: 385},

    {nversion: 21 ,correction: CorrectionL, capNum: 2232, capAlpha: 1352, capByte: 929, capKanji: 572, totalWords: 932},
    {nversion: 21 ,correction: CorrectionM, capNum: 1708, capAlpha: 1035, capByte: 711, capKanji: 438, totalWords: 714},
    {nversion: 21 ,correction: CorrectionQ, capNum: 1224, capAlpha: 742, capByte: 509, capKanji: 314, totalWords: 512},
    {nversion: 21 ,correction: CorrectionH, capNum: 969, capAlpha: 587, capByte: 403, capKanji: 248, totalWords: 406},

    {nversion: 22 ,correction: CorrectionL, capNum: 2409, capAlpha: 1460, capByte: 1003, capKanji: 618, totalWords: 1006},
    {nversion: 22 ,correction: CorrectionM, capNum: 1872, capAlpha: 1134, capByte: 779, capKanji: 480, totalWords: 782},
    {nversion: 22 ,correction: CorrectionQ, capNum: 1358, capAlpha: 823, capByte: 565, capKanji: 348, totalWords: 568},
    {nversion: 22 ,correction: CorrectionH, capNum: 1056, capAlpha: 640, capByte: 439, capKanji: 270, totalWords: 442},

    {nversion: 23 ,correction: CorrectionL, capNum: 2620, capAlpha: 1588, capByte: 1091, capKanji: 672, totalWords: 1094},
    {nversion: 23 ,correction: CorrectionM, capNum: 2059, capAlpha: 1248, capByte: 857, capKanji: 528, totalWords: 860},
    {nversion: 23 ,correction: CorrectionQ, capNum: 1468, capAlpha: 890, capByte: 611, capKanji: 376, totalWords: 614},
    {nversion: 23 ,correction: CorrectionH, capNum: 1108, capAlpha: 672, capByte: 461, capKanji: 284, totalWords: 464},

    {nversion: 24 ,correction: CorrectionL, capNum: 2812, capAlpha: 1704, capByte: 1171, capKanji: 721, totalWords: 1174},
    {nversion: 24 ,correction: CorrectionM, capNum: 2188, capAlpha: 1326, capByte: 911, capKanji: 561, totalWords: 914},
    {nversion: 24 ,correction: CorrectionQ, capNum: 1588, capAlpha: 963, capByte: 661, capKanji: 407, totalWords: 664},
    {nversion: 24 ,correction: CorrectionH, capNum: 1228, capAlpha: 744, capByte: 511, capKanji: 315, totalWords: 514},

    {nversion: 25 ,correction: CorrectionL, capNum: 3057, capAlpha: 1853, capByte: 1273, capKanji: 784, totalWords: 1276},
    {nversion: 25 ,correction: CorrectionM, capNum: 2395, capAlpha: 1451, capByte: 997, capKanji: 614, totalWords: 1000},
    {nversion: 25 ,correction: CorrectionQ, capNum: 1718, capAlpha: 1041, capByte: 715, capKanji: 440, totalWords: 718},
    {nversion: 25 ,correction: CorrectionH, capNum: 1286, capAlpha: 779, capByte: 535, capKanji: 330, totalWords: 538},

    {nversion: 26 ,correction: CorrectionL, capNum: 3283, capAlpha: 1990, capByte: 1367, capKanji: 842, totalWords: 1370},
    {nversion: 26 ,correction: CorrectionM, capNum: 2544, capAlpha: 1542, capByte: 1059, capKanji: 652, totalWords: 1062},
    {nversion: 26 ,correction: CorrectionQ, capNum: 1804, capAlpha: 1094, capByte: 751, capKanji: 462, totalWords: 754},
    {nversion: 26 ,correction: CorrectionH, capNum: 1425, capAlpha: 864, capByte: 593, capKanji: 365, totalWords: 596},

    {nversion: 27 ,correction: CorrectionL, capNum: 3517, capAlpha: 2132, capByte: 1465, capKanji: 902, totalWords: 1468},
    {nversion: 27 ,correction: CorrectionM, capNum: 2701, capAlpha: 1637, capByte: 1125, capKanji: 692, totalWords: 1128},
    {nversion: 27 ,correction: CorrectionQ, capNum: 1933, capAlpha: 1172, capByte: 805, capKanji: 496, totalWords: 808},
    {nversion: 27 ,correction: CorrectionH, capNum: 1501, capAlpha: 910, capByte: 625, capKanji: 385, totalWords: 628},

    {nversion: 28 ,correction: CorrectionL, capNum: 3669, capAlpha: 2223, capByte: 1528, capKanji: 940, totalWords: 1531},
    {nversion: 28 ,correction: CorrectionM, capNum: 2857, capAlpha: 1732, capByte: 1190, capKanji: 732, totalWords: 1193},
    {nversion: 28 ,correction: CorrectionQ, capNum: 2085, capAlpha: 1263, capByte: 868, capKanji: 534, totalWords: 871},
    {nversion: 28 ,correction: CorrectionH, capNum: 1581, capAlpha: 958, capByte: 658, capKanji: 405, totalWords: 661},

    {nversion: 29 ,correction: CorrectionL, capNum: 3909, capAlpha: 2369, capByte: 1628, capKanji: 1002, totalWords: 1631},
    {nversion: 29 ,correction: CorrectionM, capNum: 3035, capAlpha: 1839, capByte: 1264, capKanji: 778, totalWords: 1264},
    {nversion: 29 ,correction: CorrectionQ, capNum: 2181, capAlpha: 1322, capByte: 908, capKanji: 559, totalWords: 911},
    {nversion: 29 ,correction: CorrectionH, capNum: 1677, capAlpha: 1016, capByte: 698, capKanji: 430, totalWords: 701},

    {nversion: 30 ,correction: CorrectionL, capNum: 4158, capAlpha: 2520, capByte: 1732, capKanji: 1066, totalWords: 1735},
    {nversion: 30 ,correction: CorrectionM, capNum: 3289, capAlpha: 1994, capByte: 1370, capKanji: 843, totalWords: 1373},
    {nversion: 30 ,correction: CorrectionQ, capNum: 2358, capAlpha: 1429, capByte: 982, capKanji: 604, totalWords: 985},
    {nversion: 30 ,correction: CorrectionH, capNum: 1782, capAlpha: 1080, capByte: 742, capKanji: 457, totalWords: 745},

    {nversion: 31 ,correction: CorrectionL, capNum: 4417, capAlpha: 2677, capByte: 1840, capKanji: 1132, totalWords: 1843},
    {nversion: 31 ,correction: CorrectionM, capNum: 3486, capAlpha: 2113, capByte: 1452, capKanji: 894, totalWords: 1455},
    {nversion: 31 ,correction: CorrectionQ, capNum: 2473, capAlpha: 1499, capByte: 1030, capKanji: 634, totalWords: 1033},
    {nversion: 31 ,correction: CorrectionH, capNum: 1897, capAlpha: 1150, capByte: 790, capKanji: 486, totalWords: 793},

    {nversion: 32 ,correction: CorrectionL, capNum: 4686, capAlpha: 2840, capByte: 1952, capKanji: 1201, totalWords: 1955},
    {nversion: 32 ,correction: CorrectionM, capNum: 3693, capAlpha: 2238, capByte: 1538, capKanji: 947, totalWords: 1541},
    {nversion: 32 ,correction: CorrectionQ, capNum: 2670, capAlpha: 1618, capByte: 1112, capKanji: 684, totalWords: 1115},
    {nversion: 32 ,correction: CorrectionH, capNum: 2022, capAlpha: 1226, capByte: 842, capKanji: 518, totalWords: 845},

    {nversion: 33 ,correction: CorrectionL, capNum: 4965, capAlpha: 3009, capByte: 2068, capKanji: 1273, totalWords: 2071},
    {nversion: 33 ,correction: CorrectionM, capNum: 3909, capAlpha: 2369, capByte: 1628, capKanji: 1002, totalWords: 1631},
    {nversion: 33 ,correction: CorrectionQ, capNum: 2805, capAlpha: 1700, capByte: 1168, capKanji: 719, totalWords: 1171},
    {nversion: 33 ,correction: CorrectionH, capNum: 2157, capAlpha: 1307, capByte: 898, capKanji: 553, totalWords: 901},

    {nversion: 34 ,correction: CorrectionL, capNum: 5253, capAlpha: 3183, capByte: 2188, capKanji: 1347, totalWords: 2191},
    {nversion: 34 ,correction: CorrectionM, capNum: 4134, capAlpha: 2506, capByte: 1722, capKanji: 1060, totalWords: 1725},
    {nversion: 34 ,correction: CorrectionQ, capNum: 2949, capAlpha: 1787, capByte: 1228, capKanji: 756, totalWords: 1231},
    {nversion: 34 ,correction: CorrectionH, capNum: 2301, capAlpha: 1394, capByte: 958, capKanji: 590, totalWords: 961},

    {nversion: 35 ,correction: CorrectionL, capNum: 5529, capAlpha: 3351, capByte: 2303, capKanji: 1417, totalWords: 2306},
    {nversion: 35 ,correction: CorrectionM, capNum: 4343, capAlpha: 2632, capByte: 1809, capKanji: 1113, totalWords: 1812},
    {nversion: 35 ,correction: CorrectionQ, capNum: 3081, capAlpha: 1867, capByte: 1283, capKanji: 790, totalWords: 1286},
    {nversion: 35 ,correction: CorrectionH, capNum: 2361, capAlpha: 1431, capByte: 983, capKanji: 605, totalWords: 986},

    {nversion: 36 ,correction: CorrectionL, capNum: 5836, capAlpha: 3537, capByte: 2431, capKanji: 1496, totalWords: 2434},
    {nversion: 36 ,correction: CorrectionM, capNum: 4588, capAlpha: 2780, capByte: 1911, capKanji: 1176, totalWords: 1914},
    {nversion: 36 ,correction: CorrectionQ, capNum: 3244, capAlpha: 1966, capByte: 1351, capKanji: 832, totalWords: 1354},
    {nversion: 36 ,correction: CorrectionH, capNum: 2524, capAlpha: 1530, capByte: 1051, capKanji: 647, totalWords: 1054},

    {nversion: 37 ,correction: CorrectionL, capNum: 6153, capAlpha: 3729, capByte: 2563, capKanji: 1577, totalWords: 2566},
    {nversion: 37 ,correction: CorrectionM, capNum: 4775, capAlpha: 2894, capByte: 1989, capKanji: 1224, totalWords: 1992},
    {nversion: 37 ,correction: CorrectionQ, capNum: 3417, capAlpha: 2071, capByte: 1423, capKanji: 876, totalWords: 1426},
    {nversion: 37 ,correction: CorrectionH, capNum: 2625, capAlpha: 1591, capByte: 1093, capKanji: 673, totalWords: 1096},

    {nversion: 38 ,correction: CorrectionL, capNum: 6479, capAlpha: 3927, capByte: 2699, capKanji: 1661, totalWords: 2702},
    {nversion: 38 ,correction: CorrectionM, capNum: 5039, capAlpha: 3054, capByte: 2099, capKanji: 1292, totalWords: 2102},
    {nversion: 38 ,correction: CorrectionQ, capNum: 3599, capAlpha: 2181, capByte: 1499, capKanji: 923, totalWords: 1502},
    {nversion: 38 ,correction: CorrectionH, capNum: 2735, capAlpha: 1658, capByte: 1139, capKanji: 701, totalWords: 1142},

    {nversion: 39 ,correction: CorrectionL, capNum: 6743, capAlpha: 4087, capByte: 2809, capKanji: 1729, totalWords: 2812},
    {nversion: 39 ,correction: CorrectionM, capNum: 5313, capAlpha: 3220, capByte: 2213, capKanji: 1362, totalWords: 2216},
    {nversion: 39 ,correction: CorrectionQ, capNum: 3791, capAlpha: 2298, capByte: 1579, capKanji: 972, totalWords: 1582},
    {nversion: 39 ,correction: CorrectionH, capNum: 2927, capAlpha: 1774, capByte: 1219, capKanji: 750, totalWords: 1222},

    {nversion: 40 ,correction: CorrectionL, capNum: 7089, capAlpha: 4296, capByte: 2953, capKanji: 1817, totalWords: 2956},
    {nversion: 40 ,correction: CorrectionM, capNum: 5596, capAlpha: 3391, capByte: 2331, capKanji: 1435, totalWords: 2334},
    {nversion: 40 ,correction: CorrectionQ, capNum: 3993, capAlpha: 2420, capByte: 1663, capKanji: 1024, totalWords: 1666},
    {nversion: 40 ,correction: CorrectionH, capNum: 3057, capAlpha: 1852, capByte: 1273, capKanji: 784, totalWords: 1276},
  }
}

func determineVersion(input string, correction CorrectionLevel, mode EncodingMode) Version {
  versions := listVersions()
  needed := len(input)
  for _, v := range versions {
    if v.correction != correction {
      continue
    }
    capacity := 0
    switch mode {
    case Numeric:
      capacity = v.capNum
      break
    case Alphanumeric:
      capacity = v.capAlpha
      break
    case Byte:
      capacity = v.capByte
      break
    case Kanji:
      capacity = v.capKanji
      break
      
    }
    if capacity >= needed {
      return v
    }
  }
  panic("no valid version found")
}

func encode(input string, version Version, mode EncodingMode, correction CorrectionLevel) []byte {
  bytes := make([]byte, version.totalWords)

  //add mode indicator
  //left half of first byte
  var indicatorByteMask byte = 0x00;
  switch mode {
  case Numeric:
    indicatorByteMask = 0x10;
    break;
  case Alphanumeric:
    indicatorByteMask = 0x20;
    break;
  case Byte:
    indicatorByteMask = 0x40;
    break;
  case Kanji:
    indicatorByteMask = 0x80;
    break;
  }
  bytes[0] = bytes[0] | indicatorByteMask

  // add char count
  countBits := version.CharCountLength(mode)
  // fmt.Printf("need '%d' bits for the count\n", countBits)
  charCount := uint32(len(input))
  // fmt.Printf("char count '%d': '%08b'\n", charCount, u32tob(charCount))
  numBitLen := bits.Len32(charCount)
  leadZero := bits.LeadingZeros32(charCount)
  withoutLeft := charCount << leadZero
  // fmt.Printf("without left: '%08b'\n", u32tob(withoutLeft))
  neededLeft := countBits - numBitLen
  //charchountbits contains the padded number so the first n bits are the required ones
  charCountBits := withoutLeft >> uint32(neededLeft)
  // fmt.Printf("padded left: '%08b'\n", u32tob(charCountBits))
  //insert it at bit 5 (ater the mode)
  countMask := charCountBits >> 4
  maskBytes := u32tob(countMask)
  // fmt.Printf("mask: '%08b'\n", maskBytes)
  applyMask(bytes, maskBytes)

  //get encoded data
  encodedData := encodeInMode(input, mode)
  // fmt.Printf("encoded data: '%08b'\n", encodedData)

  //add data to the result
  totalOffsetBits := 4 + countBits
  bytesOffset := totalOffsetBits/8
  inByteOffset := totalOffsetBits % 8
  for i := 0; i < len(encodedData); i++ {
    bytePos := bytesOffset + i
    byteVal :=  encodedData[i]

    //apply in-byte offset
    if inByteOffset == 0 {
      bytes[bytePos] = byteVal
    } else {
      //part on this byte
      ogValue := bytes[bytePos]
      ogPart := ogValue >> (8-inByteOffset)
      ogPart = ogPart << (8-inByteOffset)
      newPart := byteVal >> inByteOffset
      byteVal = ogPart | newPart
      bytes[bytePos] = byteVal
      //part on next byte
      nextPart :=  encodedData[i] <<(8-inByteOffset)
      bytes[bytePos+1] = nextPart
    }
  }
  // fmt.Printf("after  data: %08b\n", bytes)

  //fill extra bytes
  pattern := []byte{0xEC, 0x11}
  nextByte := bytesOffset + len(encodedData)+1
  idxInsert := 0
  for nextByte < len(bytes) {
    bytes[nextByte] = pattern[idxInsert]
    idxInsert = (idxInsert + 1) % len(pattern)
    nextByte++
  }



  return bytes
}

func applyMask(val []byte, mask []byte) {
  for i := 0; i < len(val) && i < len(mask); i++ {
    val[i] |= mask[i]
  }
}

func u32tob(val uint32) []byte {
  // res := make([]byte, 4)
  // for i := uint32(0); i < 4; i++ {
  //   res[i] = byte((val >> (8*i)) & 0xff)
  // }
  // return res
  b := make([]byte, 4)
  binary.BigEndian.PutUint32(b, val)
  return b
}

func encodeInMode(input string, mode EncodingMode) []byte {
  res := []byte{}
  switch mode {
  case Numeric:
    break
  case Alphanumeric:
    return encodeAlpha(input)
  case Byte:
    break
  case Kanji:
    break
  }
  return res
}

func encodeAlpha(input string) []byte {
  table := alphaTranslator()
  bitsPerSet := 11
  numSets := len(input)/2
  oddInput := false
  if len(input) % 2 != 0 {
    numSets++
    oddInput = true
  }
  totalBits := numSets * bitsPerSet
  totalBytes := totalBits/8
  if totalBits % 8 != 0 {
    //totalBytes++
  }

  result := make([]byte, totalBytes)

  var offset uint = 0
  byteInsert := 0
  for i := 0; i < len(input)-1; i+=2 {
    c1 := rune(input[i])
    c2 := rune(input[i+1])

    n1 := table[c1]
    n2 := table[c2]
    numres := uint32((45*n1)+n2)

    //get 11bit num
    bitlen := bits.Len32(numres)
    lead := bitsPerSet - bitlen
    withoutLead := numres << (bits.LeadingZeros32(numres)-lead)
    //11bits and then 0s

    //add padding for in-byte offset
    padded := withoutLead >> offset
    numBytes := u32tob(padded)
    byteInsert = (bitsPerSet * (i/2)) / 8
    for i := 0; i < len(numBytes) && byteInsert+i < len(result); i++ {
      result[byteInsert+i] |= numBytes[i]
    }
    offset = (offset + uint(bitsPerSet)) % 8

  }

  if oddInput {
    c := rune(input[len(input)-1])
    numres := uint32(table[c])

    //get 11bit num
    bitlen := bits.Len32(numres)
    lead := 6 - bitlen
    withoutLead := numres << (bits.LeadingZeros32(numres)-lead)
    //11bits and then 0s

    //add padding for in-byte offset
    padded := withoutLead >> offset
    numBytes := u32tob(padded)

    bitsInserted := bitsPerSet * ((len(input)/2))
    bytesUsed := bitsInserted/8

    for i := 0; i < len(numBytes) && bytesUsed+i < len(result) && i<2; i++ {
      result[bytesUsed+i] |= numBytes[i]
    }

  }


  return result
}

func alphaTranslator() map[rune]int {
  return map[rune]int {
    '0': 0,
    '1': 1,
    '2': 2,
    '3': 3,
    '4': 4,
    '5': 5,
    '6': 6,
    '7': 7,
    '8': 8,
    '9': 9,
    'A': 10,
    'B': 11,
    'C': 12,
    'D': 13,
    'E': 14,
    'F': 15,
    'G': 16,
    'H': 17,
    'I': 18,
    'J': 19,
    'K': 20,
    'L': 21,
    'M': 22,
    'N': 23,
    'O': 24,
    'P': 25,
    'Q': 26,
    'R': 27,
    'S': 28,
    'T': 29,
    'U': 30,
    'V': 31,
    'W': 32,
    'X': 33,
    'Y': 34,
    'Z': 35,
    ' ': 36,
    '$': 37,
    '%': 38,
    '*': 39,
    '+': 40,
    '-': 41,
    '.': 42,
    '/': 43,
    ':': 44,
  }
}




