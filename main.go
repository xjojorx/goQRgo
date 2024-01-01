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
  blocksGroup1 int
  wordsBlockGroup1 int
  blocksGroup2 int
  wordsBlockGroup2 int
  ecWordsBlock int
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
  corrLvl := CorrectionM //should read from args

  mode := encodingFormat(input)
  fmt.Printf("mode: %#v\n", mode)

  version := determineVersion(input, corrLvl, mode)

  fmt.Printf("input: '%s', mode: '%d', correction: '%s', version: '%d'\n", input, mode, string(corrLvl), version.nversion)

  encoded := encode(input, version, mode, corrLvl)

  fmt.Printf("encoded as: %08b\n", encoded)

  errorCorrection(encoded, version)

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
    {nversion: 1, correction: CorrectionL, capNum: 41, capAlpha: 25, capByte: 17, capKanji: 10, totalWords: 19, blocksGroup1: 1, wordsBlockGroup1: 19, blocksGroup2: 0, wordsBlockGroup2: 0, ecWordsBlock: 7},
    {nversion: 1, correction: CorrectionM, capNum: 34, capAlpha: 20, capByte: 14, capKanji: 8, totalWords: 16, blocksGroup1: 1, wordsBlockGroup1: 16, blocksGroup2: 0, wordsBlockGroup2: 0, ecWordsBlock: 10},
    {nversion: 1, correction: CorrectionQ, capNum: 27, capAlpha: 16, capByte: 11, capKanji: 7, totalWords: 13, blocksGroup1: 1, wordsBlockGroup1: 13, blocksGroup2: 0, wordsBlockGroup2: 0, ecWordsBlock: 13},
    {nversion: 1, correction: CorrectionH, capNum: 17, capAlpha: 10, capByte: 7, capKanji: 4, totalWords: 9, blocksGroup1: 1, wordsBlockGroup1: 9, blocksGroup2: 0, wordsBlockGroup2: 0, ecWordsBlock: 17},

    {nversion: 2 ,correction: CorrectionL, capNum: 77, capAlpha: 47, capByte: 32, capKanji: 20, totalWords: 34, blocksGroup1: 1, wordsBlockGroup1: 34, blocksGroup2: 0, wordsBlockGroup2: 0, ecWordsBlock: 10},
    {nversion: 2 ,correction: CorrectionM, capNum: 63, capAlpha: 38, capByte: 26, capKanji: 16, totalWords: 28, blocksGroup1: 1, wordsBlockGroup1: 28, blocksGroup2: 0, wordsBlockGroup2: 0, ecWordsBlock: 16},
    {nversion: 2 ,correction: CorrectionQ, capNum: 48, capAlpha: 29, capByte: 20, capKanji: 12, totalWords: 22, blocksGroup1: 1, wordsBlockGroup1: 22, blocksGroup2: 0, wordsBlockGroup2: 0, ecWordsBlock: 22},
    {nversion: 2 ,correction: CorrectionH, capNum: 34, capAlpha: 20, capByte: 14, capKanji: 8, totalWords: 16, blocksGroup1: 1, wordsBlockGroup1: 16, blocksGroup2: 0, wordsBlockGroup2: 0, ecWordsBlock: 28},

    {nversion: 3 ,correction: CorrectionL, capNum: 127, capAlpha: 77, capByte: 53, capKanji: 32, totalWords: 55, blocksGroup1: 1, wordsBlockGroup1: 55, blocksGroup2: 0, wordsBlockGroup2: 0, ecWordsBlock: 15},
    {nversion: 3 ,correction: CorrectionM, capNum: 101, capAlpha: 61, capByte: 42, capKanji: 26, totalWords: 44, blocksGroup1: 1, wordsBlockGroup1: 44, blocksGroup2: 0, wordsBlockGroup2: 0, ecWordsBlock: 26},
    {nversion: 3 ,correction: CorrectionQ, capNum: 77, capAlpha: 47, capByte: 32, capKanji: 20, totalWords: 34, blocksGroup1: 2, wordsBlockGroup1: 17, blocksGroup2: 0, wordsBlockGroup2: 0, ecWordsBlock: 18},
    {nversion: 3 ,correction: CorrectionH, capNum: 58, capAlpha: 35, capByte: 24, capKanji: 15, totalWords: 26, blocksGroup1: 2, wordsBlockGroup1: 13, blocksGroup2: 0, wordsBlockGroup2: 0, ecWordsBlock: 22},

    {nversion: 4 ,correction: CorrectionL, capNum: 187, capAlpha: 114, capByte: 78, capKanji: 48, totalWords: 80, blocksGroup1: 1, wordsBlockGroup1: 80, blocksGroup2: 0, wordsBlockGroup2: 0, ecWordsBlock: 20},
    {nversion: 4 ,correction: CorrectionM, capNum: 149, capAlpha: 90, capByte: 62, capKanji: 38, totalWords: 64, blocksGroup1: 2, wordsBlockGroup1: 32, blocksGroup2: 0, wordsBlockGroup2: 0, ecWordsBlock: 18},
    {nversion: 4 ,correction: CorrectionQ, capNum: 111, capAlpha: 67, capByte: 46, capKanji: 28, totalWords: 48, blocksGroup1: 2, wordsBlockGroup1: 24, blocksGroup2: 0, wordsBlockGroup2: 0, ecWordsBlock: 26},
    {nversion: 4 ,correction: CorrectionH, capNum: 82, capAlpha: 50, capByte: 34, capKanji: 21, totalWords: 36, blocksGroup1: 4, wordsBlockGroup1: 9, blocksGroup2: 0, wordsBlockGroup2: 0, ecWordsBlock: 16},

    {nversion: 5 ,correction: CorrectionL, capNum: 255, capAlpha: 154, capByte: 106, capKanji: 65, totalWords: 108, blocksGroup1: 1, wordsBlockGroup1: 108, blocksGroup2: 0, wordsBlockGroup2: 0, ecWordsBlock: 26},
    {nversion: 5 ,correction: CorrectionM, capNum: 202, capAlpha: 122, capByte: 84, capKanji: 52, totalWords: 86, blocksGroup1: 2, wordsBlockGroup1: 43, blocksGroup2: 0, wordsBlockGroup2: 0, ecWordsBlock: 24},
    {nversion: 5 ,correction: CorrectionQ, capNum: 144, capAlpha: 87, capByte: 60, capKanji: 37, totalWords: 62, blocksGroup1: 2, wordsBlockGroup1: 15, blocksGroup2: 2, wordsBlockGroup2: 16 , ecWordsBlock: 18},
    {nversion: 5 ,correction: CorrectionH, capNum: 106, capAlpha: 64, capByte: 44, capKanji: 27, totalWords: 46, blocksGroup1: 2, wordsBlockGroup1: 11, blocksGroup2: 2, wordsBlockGroup2: 12 , ecWordsBlock: 22},

    {nversion: 6 ,correction: CorrectionL, capNum: 322, capAlpha: 195, capByte: 134, capKanji: 82, totalWords: 136, blocksGroup1: 2, wordsBlockGroup1: 68, blocksGroup2: 0, wordsBlockGroup2: 0, ecWordsBlock: 18},
    {nversion: 6 ,correction: CorrectionM, capNum: 255, capAlpha: 154, capByte: 106, capKanji: 65, totalWords: 108, blocksGroup1: 4, wordsBlockGroup1: 27, blocksGroup2: 0, wordsBlockGroup2: 0, ecWordsBlock: 16},
    {nversion: 6 ,correction: CorrectionQ, capNum: 178, capAlpha: 108, capByte: 74, capKanji: 45, totalWords: 76, blocksGroup1: 4, wordsBlockGroup1: 19, blocksGroup2: 0, wordsBlockGroup2: 0, ecWordsBlock: 24},
    {nversion: 6 ,correction: CorrectionH, capNum: 139, capAlpha: 84, capByte: 58, capKanji: 36, totalWords: 60, blocksGroup1: 4, wordsBlockGroup1: 15, blocksGroup2: 0, wordsBlockGroup2: 0, ecWordsBlock: 28},

    {nversion: 7 ,correction: CorrectionL, capNum: 370, capAlpha: 224, capByte: 154, capKanji: 95, totalWords: 156, blocksGroup1: 2, wordsBlockGroup1: 78, blocksGroup2: 0, wordsBlockGroup2: 0, ecWordsBlock: 20},
    {nversion: 7 ,correction: CorrectionM, capNum: 293, capAlpha: 178, capByte: 122, capKanji: 75, totalWords: 124, blocksGroup1: 4, wordsBlockGroup1: 31, blocksGroup2: 0, wordsBlockGroup2: 0, ecWordsBlock: 18},
    {nversion: 7 ,correction: CorrectionQ, capNum: 207, capAlpha: 125, capByte: 86, capKanji: 53, totalWords: 88, blocksGroup1: 2, wordsBlockGroup1: 14, blocksGroup2: 4, wordsBlockGroup2: 15 , ecWordsBlock: 18},
    {nversion: 7 ,correction: CorrectionH, capNum: 154, capAlpha: 93, capByte: 64, capKanji: 39, totalWords: 66, blocksGroup1: 4, wordsBlockGroup1: 13, blocksGroup2: 1, wordsBlockGroup2: 14 , ecWordsBlock: 26},

    {nversion: 8 ,correction: CorrectionL, capNum: 461, capAlpha: 279, capByte: 192, capKanji: 118, totalWords: 194, blocksGroup1: 2, wordsBlockGroup1: 97, blocksGroup2: 0, wordsBlockGroup2: 0, ecWordsBlock: 24},
    {nversion: 8 ,correction: CorrectionM, capNum: 365, capAlpha: 221, capByte: 152, capKanji: 93, totalWords: 154, blocksGroup1: 2, wordsBlockGroup1: 38, blocksGroup2: 2, wordsBlockGroup2: 39 , ecWordsBlock: 22},
    {nversion: 8 ,correction: CorrectionQ, capNum: 259, capAlpha: 157, capByte: 108, capKanji: 66, totalWords: 110, blocksGroup1: 4, wordsBlockGroup1: 18, blocksGroup2: 2, wordsBlockGroup2: 19 , ecWordsBlock: 22},
    {nversion: 8 ,correction: CorrectionH, capNum: 202, capAlpha: 122, capByte: 84, capKanji: 52, totalWords: 86, blocksGroup1: 4, wordsBlockGroup1: 14, blocksGroup2: 2, wordsBlockGroup2: 15 , ecWordsBlock: 26},

    {nversion: 9 ,correction: CorrectionL, capNum: 552, capAlpha: 335, capByte: 230, capKanji: 141, totalWords: 232, blocksGroup1: 2, wordsBlockGroup1: 116, blocksGroup2: 0, wordsBlockGroup2: 0, ecWordsBlock: 30},
    {nversion: 9 ,correction: CorrectionM, capNum: 432, capAlpha: 262, capByte: 180, capKanji: 111, totalWords: 182, blocksGroup1: 3, wordsBlockGroup1: 36, blocksGroup2: 2, wordsBlockGroup2: 37 , ecWordsBlock: 22},
    {nversion: 9 ,correction: CorrectionQ, capNum: 312, capAlpha: 189, capByte: 130, capKanji: 80, totalWords: 132, blocksGroup1: 4, wordsBlockGroup1: 16, blocksGroup2: 4, wordsBlockGroup2: 17 , ecWordsBlock: 20},
    {nversion: 9 ,correction: CorrectionH, capNum: 235, capAlpha: 143, capByte: 98, capKanji: 60, totalWords: 100, blocksGroup1: 4, wordsBlockGroup1: 12, blocksGroup2: 4, wordsBlockGroup2: 13 , ecWordsBlock: 24},

    {nversion: 10 ,correction: CorrectionL, capNum: 652, capAlpha: 395, capByte: 271, capKanji: 167, totalWords: 274, blocksGroup1: 2, wordsBlockGroup1: 68, blocksGroup2: 2, wordsBlockGroup2: 69 , ecWordsBlock: 18},
    {nversion: 10 ,correction: CorrectionM, capNum: 513, capAlpha: 311, capByte: 213, capKanji: 131, totalWords: 216, blocksGroup1: 4, wordsBlockGroup1: 43, blocksGroup2: 1, wordsBlockGroup2: 44 , ecWordsBlock: 26},
    {nversion: 10 ,correction: CorrectionQ, capNum: 364, capAlpha: 221, capByte: 151, capKanji: 93, totalWords: 154, blocksGroup1: 6, wordsBlockGroup1: 19, blocksGroup2: 2, wordsBlockGroup2: 20 , ecWordsBlock: 24},
    {nversion: 10 ,correction: CorrectionH, capNum: 288, capAlpha: 174, capByte: 119, capKanji: 74, totalWords: 122, blocksGroup1: 6, wordsBlockGroup1: 15, blocksGroup2: 2, wordsBlockGroup2: 16 , ecWordsBlock: 28},

    {nversion: 11 ,correction: CorrectionL, capNum: 772, capAlpha: 468, capByte: 321, capKanji: 198, totalWords: 324, blocksGroup1: 4, wordsBlockGroup1: 81, blocksGroup2: 0, wordsBlockGroup2: 0, ecWordsBlock: 20},
    {nversion: 11 ,correction: CorrectionM, capNum: 604, capAlpha: 366, capByte: 251, capKanji: 155, totalWords: 254, blocksGroup1: 1, wordsBlockGroup1: 50, blocksGroup2: 4, wordsBlockGroup2: 51 , ecWordsBlock: 30},
    {nversion: 11 ,correction: CorrectionQ, capNum: 427, capAlpha: 259, capByte: 177, capKanji: 109, totalWords: 180, blocksGroup1: 4, wordsBlockGroup1: 22, blocksGroup2: 4, wordsBlockGroup2: 23 , ecWordsBlock: 28},
    {nversion: 11 ,correction: CorrectionH, capNum: 331, capAlpha: 200, capByte: 137, capKanji: 85, totalWords: 140, blocksGroup1: 3, wordsBlockGroup1: 12, blocksGroup2: 8, wordsBlockGroup2: 13 , ecWordsBlock: 24},

    {nversion: 12 ,correction: CorrectionL, capNum: 883, capAlpha: 535, capByte: 367, capKanji: 226, totalWords: 370, blocksGroup1: 2, wordsBlockGroup1: 92, blocksGroup2: 2, wordsBlockGroup2: 93 , ecWordsBlock: 24},
    {nversion: 12 ,correction: CorrectionM, capNum: 691, capAlpha: 419, capByte: 287, capKanji: 177, totalWords: 290, blocksGroup1: 6, wordsBlockGroup1: 36, blocksGroup2: 2, wordsBlockGroup2: 37 , ecWordsBlock: 22},
    {nversion: 12 ,correction: CorrectionQ, capNum: 489, capAlpha: 296, capByte: 203, capKanji: 125, totalWords: 206, blocksGroup1: 4, wordsBlockGroup1: 20, blocksGroup2: 6, wordsBlockGroup2: 21 , ecWordsBlock: 26},
    {nversion: 12 ,correction: CorrectionH, capNum: 374, capAlpha: 227, capByte: 155, capKanji: 96, totalWords: 158, blocksGroup1: 7, wordsBlockGroup1: 14, blocksGroup2: 4, wordsBlockGroup2: 15 , ecWordsBlock: 28},

    {nversion: 13 ,correction: CorrectionL, capNum: 1022, capAlpha: 619, capByte: 425, capKanji: 262, totalWords: 428, blocksGroup1: 4, wordsBlockGroup1: 107, blocksGroup2: 0, wordsBlockGroup2: 0, ecWordsBlock: 26},
    {nversion: 13 ,correction: CorrectionM, capNum: 796, capAlpha: 483, capByte: 331, capKanji: 204, totalWords: 334, blocksGroup1: 8, wordsBlockGroup1: 37, blocksGroup2: 1, wordsBlockGroup2: 38 , ecWordsBlock: 22},
    {nversion: 13 ,correction: CorrectionQ, capNum: 580, capAlpha: 352, capByte: 241, capKanji: 149, totalWords: 244, blocksGroup1: 8, wordsBlockGroup1: 20, blocksGroup2: 4, wordsBlockGroup2: 21 , ecWordsBlock: 24},
    {nversion: 13 ,correction: CorrectionH, capNum: 427, capAlpha: 259, capByte: 177, capKanji: 109, totalWords: 180, blocksGroup1: 12, wordsBlockGroup1: 11, blocksGroup2: 4, wordsBlockGroup2: 12 , ecWordsBlock: 22},

    {nversion: 14 ,correction: CorrectionL, capNum: 1101, capAlpha: 667, capByte: 458, capKanji: 282, totalWords: 461, blocksGroup1: 3, wordsBlockGroup1: 115, blocksGroup2: 1, wordsBlockGroup2: 116 , ecWordsBlock: 30},
    {nversion: 14 ,correction: CorrectionM, capNum: 871, capAlpha: 528, capByte: 362, capKanji: 223, totalWords: 365, blocksGroup1: 4, wordsBlockGroup1: 40, blocksGroup2: 5, wordsBlockGroup2: 41 , ecWordsBlock: 24},
    {nversion: 14 ,correction: CorrectionQ, capNum: 621, capAlpha: 376, capByte: 258, capKanji: 159, totalWords: 261, blocksGroup1: 11, wordsBlockGroup1: 16, blocksGroup2: 5, wordsBlockGroup2: 17 , ecWordsBlock: 20},
    {nversion: 14 ,correction: CorrectionH, capNum: 468, capAlpha: 283, capByte: 194, capKanji: 120, totalWords: 197, blocksGroup1: 11, wordsBlockGroup1: 12, blocksGroup2: 5, wordsBlockGroup2: 13 , ecWordsBlock: 24},

    {nversion: 15 ,correction: CorrectionL, capNum: 1250, capAlpha: 758, capByte: 520, capKanji: 320, totalWords: 523, blocksGroup1: 5, wordsBlockGroup1: 87, blocksGroup2: 1, wordsBlockGroup2: 88 , ecWordsBlock: 22},
    {nversion: 15 ,correction: CorrectionM, capNum: 991, capAlpha: 600, capByte: 412, capKanji: 254, totalWords: 415, blocksGroup1: 5, wordsBlockGroup1: 41, blocksGroup2: 5, wordsBlockGroup2: 42 , ecWordsBlock: 24},
    {nversion: 15 ,correction: CorrectionQ, capNum: 703, capAlpha: 426, capByte: 292, capKanji: 180, totalWords: 295, blocksGroup1: 5, wordsBlockGroup1: 24, blocksGroup2: 7, wordsBlockGroup2: 25 , ecWordsBlock: 30},
    {nversion: 15 ,correction: CorrectionH, capNum: 530, capAlpha: 321, capByte: 220, capKanji: 136, totalWords: 223, blocksGroup1: 11, wordsBlockGroup1: 12, blocksGroup2: 7, wordsBlockGroup2: 13 , ecWordsBlock: 24},

    {nversion: 16 ,correction: CorrectionL, capNum: 1408, capAlpha: 854, capByte: 586, capKanji: 361, totalWords: 589, blocksGroup1: 5, wordsBlockGroup1: 98, blocksGroup2: 1, wordsBlockGroup2: 99 , ecWordsBlock: 24},
    {nversion: 16 ,correction: CorrectionM, capNum: 1082, capAlpha: 656, capByte: 450, capKanji: 277, totalWords: 453, blocksGroup1: 7, wordsBlockGroup1: 45, blocksGroup2: 3, wordsBlockGroup2: 46 , ecWordsBlock: 28},
    {nversion: 16 ,correction: CorrectionQ, capNum: 775, capAlpha: 470, capByte: 322, capKanji: 198, totalWords: 325, blocksGroup1: 15, wordsBlockGroup1: 19, blocksGroup2: 2, wordsBlockGroup2: 20 , ecWordsBlock: 24},
    {nversion: 16 ,correction: CorrectionH, capNum: 602, capAlpha: 365, capByte: 250, capKanji: 154, totalWords: 253, blocksGroup1: 3, wordsBlockGroup1: 15, blocksGroup2: 13, wordsBlockGroup2: 16 , ecWordsBlock: 30},

    {nversion: 17 ,correction: CorrectionL, capNum: 1548, capAlpha: 938, capByte: 644, capKanji: 397, totalWords: 647, blocksGroup1: 1, wordsBlockGroup1: 107, blocksGroup2: 5, wordsBlockGroup2: 108 , ecWordsBlock: 28},
    {nversion: 17 ,correction: CorrectionM, capNum: 1212, capAlpha: 734, capByte: 504, capKanji: 310, totalWords: 507, blocksGroup1: 10, wordsBlockGroup1: 46, blocksGroup2: 1, wordsBlockGroup2: 47 , ecWordsBlock: 28},
    {nversion: 17 ,correction: CorrectionQ, capNum: 876, capAlpha: 531, capByte: 364, capKanji: 224, totalWords: 367, blocksGroup1: 1, wordsBlockGroup1: 22, blocksGroup2: 15, wordsBlockGroup2: 23 , ecWordsBlock: 28},
    {nversion: 17 ,correction: CorrectionH, capNum: 674, capAlpha: 408, capByte: 280, capKanji: 173, totalWords: 283, blocksGroup1: 2, wordsBlockGroup1: 14, blocksGroup2: 17, wordsBlockGroup2: 15 , ecWordsBlock: 28},

    {nversion: 18 ,correction: CorrectionL, capNum: 1725, capAlpha: 1046, capByte: 718, capKanji: 442, totalWords: 721, blocksGroup1: 5, wordsBlockGroup1: 120, blocksGroup2: 1, wordsBlockGroup2: 121 , ecWordsBlock: 30},
    {nversion: 18 ,correction: CorrectionM, capNum: 1346, capAlpha: 816, capByte: 560, capKanji: 345, totalWords: 563, blocksGroup1: 9, wordsBlockGroup1: 43, blocksGroup2: 4, wordsBlockGroup2: 44 , ecWordsBlock: 26},
    {nversion: 18 ,correction: CorrectionQ, capNum: 948, capAlpha: 574, capByte: 394, capKanji: 243, totalWords: 397, blocksGroup1: 17, wordsBlockGroup1: 22, blocksGroup2: 1, wordsBlockGroup2: 23 , ecWordsBlock: 28},
    {nversion: 18 ,correction: CorrectionH, capNum: 746, capAlpha: 452, capByte: 310, capKanji: 191, totalWords: 313, blocksGroup1: 2, wordsBlockGroup1: 14, blocksGroup2: 19, wordsBlockGroup2: 15 , ecWordsBlock: 28},

    {nversion: 19 ,correction: CorrectionL, capNum: 1903, capAlpha: 1153, capByte: 792, capKanji: 488, totalWords: 795, blocksGroup1: 3, wordsBlockGroup1: 113, blocksGroup2: 4, wordsBlockGroup2: 114 , ecWordsBlock: 28},
    {nversion: 19 ,correction: CorrectionM, capNum: 1500, capAlpha: 909, capByte: 624, capKanji: 384, totalWords: 627, blocksGroup1: 3, wordsBlockGroup1: 44, blocksGroup2: 11, wordsBlockGroup2: 45 , ecWordsBlock: 26},
    {nversion: 19 ,correction: CorrectionQ, capNum: 1063, capAlpha: 644, capByte: 442, capKanji: 272, totalWords: 445, blocksGroup1: 17, wordsBlockGroup1: 21, blocksGroup2: 4, wordsBlockGroup2: 22 , ecWordsBlock: 26},
    {nversion: 19 ,correction: CorrectionH, capNum: 813, capAlpha: 493, capByte: 338, capKanji: 208, totalWords: 341, blocksGroup1: 9, wordsBlockGroup1: 13, blocksGroup2: 16, wordsBlockGroup2: 14 , ecWordsBlock: 26},

    {nversion: 20 ,correction: CorrectionL, capNum: 2061, capAlpha: 1249, capByte: 858, capKanji: 528, totalWords: 861, blocksGroup1: 3, wordsBlockGroup1: 107, blocksGroup2: 5, wordsBlockGroup2: 108 , ecWordsBlock: 28},
    {nversion: 20 ,correction: CorrectionM, capNum: 1600, capAlpha: 970, capByte: 666, capKanji: 410, totalWords: 669, blocksGroup1: 3, wordsBlockGroup1: 41, blocksGroup2: 13, wordsBlockGroup2: 42 , ecWordsBlock: 26},
    {nversion: 20 ,correction: CorrectionQ, capNum: 1159, capAlpha: 702, capByte: 482, capKanji: 297, totalWords: 485, blocksGroup1: 15, wordsBlockGroup1: 24, blocksGroup2: 5, wordsBlockGroup2: 25 , ecWordsBlock: 30},
    {nversion: 20 ,correction: CorrectionH, capNum: 919, capAlpha: 557, capByte: 382, capKanji: 235, totalWords: 385, blocksGroup1: 15, wordsBlockGroup1: 15, blocksGroup2: 10, wordsBlockGroup2: 16 , ecWordsBlock: 28},

    {nversion: 21 ,correction: CorrectionL, capNum: 2232, capAlpha: 1352, capByte: 929, capKanji: 572, totalWords: 932, blocksGroup1: 4, wordsBlockGroup1: 116, blocksGroup2: 4, wordsBlockGroup2: 117 , ecWordsBlock: 28},
    {nversion: 21 ,correction: CorrectionM, capNum: 1708, capAlpha: 1035, capByte: 711, capKanji: 438, totalWords: 714, blocksGroup1: 17, wordsBlockGroup1: 42, blocksGroup2: 0, wordsBlockGroup2: 0, ecWordsBlock: 26},
    {nversion: 21 ,correction: CorrectionQ, capNum: 1224, capAlpha: 742, capByte: 509, capKanji: 314, totalWords: 512, blocksGroup1: 17, wordsBlockGroup1: 22, blocksGroup2: 6, wordsBlockGroup2: 23 , ecWordsBlock: 28},
    {nversion: 21 ,correction: CorrectionH, capNum: 969, capAlpha: 587, capByte: 403, capKanji: 248, totalWords: 406, blocksGroup1: 19, wordsBlockGroup1: 16, blocksGroup2: 6, wordsBlockGroup2: 17 , ecWordsBlock: 30},

    {nversion: 22 ,correction: CorrectionL, capNum: 2409, capAlpha: 1460, capByte: 1003, capKanji: 618, totalWords: 1006, blocksGroup1: 2, wordsBlockGroup1: 111, blocksGroup2: 7, wordsBlockGroup2: 112 , ecWordsBlock: 28},
    {nversion: 22 ,correction: CorrectionM, capNum: 1872, capAlpha: 1134, capByte: 779, capKanji: 480, totalWords: 782, blocksGroup1: 17, wordsBlockGroup1: 46, blocksGroup2: 0, wordsBlockGroup2: 0, ecWordsBlock: 28},
    {nversion: 22 ,correction: CorrectionQ, capNum: 1358, capAlpha: 823, capByte: 565, capKanji: 348, totalWords: 568, blocksGroup1: 7, wordsBlockGroup1: 24, blocksGroup2: 16, wordsBlockGroup2: 25 , ecWordsBlock: 30},
    {nversion: 22 ,correction: CorrectionH, capNum: 1056, capAlpha: 640, capByte: 439, capKanji: 270, totalWords: 442, blocksGroup1: 34, wordsBlockGroup1: 13, blocksGroup2: 0, wordsBlockGroup2: 0, ecWordsBlock: 24},

    {nversion: 23 ,correction: CorrectionL, capNum: 2620, capAlpha: 1588, capByte: 1091, capKanji: 672, totalWords: 1094, blocksGroup1: 4, wordsBlockGroup1: 121, blocksGroup2: 5, wordsBlockGroup2: 122 , ecWordsBlock: 30},
    {nversion: 23 ,correction: CorrectionM, capNum: 2059, capAlpha: 1248, capByte: 857, capKanji: 528, totalWords: 860, blocksGroup1: 4, wordsBlockGroup1: 47, blocksGroup2: 14, wordsBlockGroup2: 48 , ecWordsBlock: 28},
    {nversion: 23 ,correction: CorrectionQ, capNum: 1468, capAlpha: 890, capByte: 611, capKanji: 376, totalWords: 614, blocksGroup1: 11, wordsBlockGroup1: 24, blocksGroup2: 14, wordsBlockGroup2: 25 , ecWordsBlock: 30},
    {nversion: 23 ,correction: CorrectionH, capNum: 1108, capAlpha: 672, capByte: 461, capKanji: 284, totalWords: 464, blocksGroup1: 16, wordsBlockGroup1: 15, blocksGroup2: 14, wordsBlockGroup2: 16 , ecWordsBlock: 30},

    {nversion: 24 ,correction: CorrectionL, capNum: 2812, capAlpha: 1704, capByte: 1171, capKanji: 721, totalWords: 1174, blocksGroup1: 6, wordsBlockGroup1: 117, blocksGroup2: 4, wordsBlockGroup2: 118 , ecWordsBlock: 30},
    {nversion: 24 ,correction: CorrectionM, capNum: 2188, capAlpha: 1326, capByte: 911, capKanji: 561, totalWords: 914, blocksGroup1: 6, wordsBlockGroup1: 45, blocksGroup2: 14, wordsBlockGroup2: 46 , ecWordsBlock: 28},
    {nversion: 24 ,correction: CorrectionQ, capNum: 1588, capAlpha: 963, capByte: 661, capKanji: 407, totalWords: 664, blocksGroup1: 11, wordsBlockGroup1: 24, blocksGroup2: 16, wordsBlockGroup2: 25 , ecWordsBlock: 30},
    {nversion: 24 ,correction: CorrectionH, capNum: 1228, capAlpha: 744, capByte: 511, capKanji: 315, totalWords: 514, blocksGroup1: 30, wordsBlockGroup1: 16, blocksGroup2: 2, wordsBlockGroup2: 17 , ecWordsBlock: 30},

    {nversion: 25 ,correction: CorrectionL, capNum: 3057, capAlpha: 1853, capByte: 1273, capKanji: 784, totalWords: 1276, blocksGroup1: 8, wordsBlockGroup1: 106, blocksGroup2: 4, wordsBlockGroup2: 107 , ecWordsBlock: 26},
    {nversion: 25 ,correction: CorrectionM, capNum: 2395, capAlpha: 1451, capByte: 997, capKanji: 614, totalWords: 1000, blocksGroup1: 8, wordsBlockGroup1: 47, blocksGroup2: 13, wordsBlockGroup2: 48 , ecWordsBlock: 28},
    {nversion: 25 ,correction: CorrectionQ, capNum: 1718, capAlpha: 1041, capByte: 715, capKanji: 440, totalWords: 718, blocksGroup1: 7, wordsBlockGroup1: 24, blocksGroup2: 22, wordsBlockGroup2: 25 , ecWordsBlock: 30},
    {nversion: 25 ,correction: CorrectionH, capNum: 1286, capAlpha: 779, capByte: 535, capKanji: 330, totalWords: 538, blocksGroup1: 22, wordsBlockGroup1: 15, blocksGroup2: 13, wordsBlockGroup2: 16 , ecWordsBlock: 30},

    {nversion: 26 ,correction: CorrectionL, capNum: 3283, capAlpha: 1990, capByte: 1367, capKanji: 842, totalWords: 1370, blocksGroup1: 10, wordsBlockGroup1: 114, blocksGroup2: 2, wordsBlockGroup2: 115 , ecWordsBlock: 28},
    {nversion: 26 ,correction: CorrectionM, capNum: 2544, capAlpha: 1542, capByte: 1059, capKanji: 652, totalWords: 1062, blocksGroup1: 19, wordsBlockGroup1: 46, blocksGroup2: 4, wordsBlockGroup2: 47 , ecWordsBlock: 28},
    {nversion: 26 ,correction: CorrectionQ, capNum: 1804, capAlpha: 1094, capByte: 751, capKanji: 462, totalWords: 754, blocksGroup1: 28, wordsBlockGroup1: 22, blocksGroup2: 6, wordsBlockGroup2: 23 , ecWordsBlock: 28},
    {nversion: 26 ,correction: CorrectionH, capNum: 1425, capAlpha: 864, capByte: 593, capKanji: 365, totalWords: 596, blocksGroup1: 33, wordsBlockGroup1: 16, blocksGroup2: 4, wordsBlockGroup2: 17 , ecWordsBlock: 30},

    {nversion: 27 ,correction: CorrectionL, capNum: 3517, capAlpha: 2132, capByte: 1465, capKanji: 902, totalWords: 1468, blocksGroup1: 8, wordsBlockGroup1: 122, blocksGroup2: 4, wordsBlockGroup2: 123 , ecWordsBlock: 30},
    {nversion: 27 ,correction: CorrectionM, capNum: 2701, capAlpha: 1637, capByte: 1125, capKanji: 692, totalWords: 1128, blocksGroup1: 22, wordsBlockGroup1: 45, blocksGroup2: 3, wordsBlockGroup2: 46 , ecWordsBlock: 28},
    {nversion: 27 ,correction: CorrectionQ, capNum: 1933, capAlpha: 1172, capByte: 805, capKanji: 496, totalWords: 808, blocksGroup1: 8, wordsBlockGroup1: 23, blocksGroup2: 26, wordsBlockGroup2: 24 , ecWordsBlock: 30},
    {nversion: 27 ,correction: CorrectionH, capNum: 1501, capAlpha: 910, capByte: 625, capKanji: 385, totalWords: 628, blocksGroup1: 12, wordsBlockGroup1: 15, blocksGroup2: 28, wordsBlockGroup2: 16 , ecWordsBlock: 30},

    {nversion: 28 ,correction: CorrectionL, capNum: 3669, capAlpha: 2223, capByte: 1528, capKanji: 940, totalWords: 1531, blocksGroup1: 3, wordsBlockGroup1: 117, blocksGroup2: 10, wordsBlockGroup2: 118 , ecWordsBlock: 30},
    {nversion: 28 ,correction: CorrectionM, capNum: 2857, capAlpha: 1732, capByte: 1190, capKanji: 732, totalWords: 1193, blocksGroup1: 3, wordsBlockGroup1: 45, blocksGroup2: 23, wordsBlockGroup2: 46 , ecWordsBlock: 28},
    {nversion: 28 ,correction: CorrectionQ, capNum: 2085, capAlpha: 1263, capByte: 868, capKanji: 534, totalWords: 871, blocksGroup1: 4, wordsBlockGroup1: 24, blocksGroup2: 31, wordsBlockGroup2: 25 , ecWordsBlock: 30},
    {nversion: 28 ,correction: CorrectionH, capNum: 1581, capAlpha: 958, capByte: 658, capKanji: 405, totalWords: 661, blocksGroup1: 11, wordsBlockGroup1: 15, blocksGroup2: 31, wordsBlockGroup2: 16 , ecWordsBlock: 30},

    {nversion: 29 ,correction: CorrectionL, capNum: 3909, capAlpha: 2369, capByte: 1628, capKanji: 1002, totalWords: 1631, blocksGroup1: 7, wordsBlockGroup1: 116, blocksGroup2: 7, wordsBlockGroup2: 117 , ecWordsBlock: 30},
    {nversion: 29 ,correction: CorrectionM, capNum: 3035, capAlpha: 1839, capByte: 1264, capKanji: 778, totalWords: 1264, blocksGroup1: 21, wordsBlockGroup1: 45, blocksGroup2: 7, wordsBlockGroup2: 46 , ecWordsBlock: 28},
    {nversion: 29 ,correction: CorrectionQ, capNum: 2181, capAlpha: 1322, capByte: 908, capKanji: 559, totalWords: 911, blocksGroup1: 1, wordsBlockGroup1: 23, blocksGroup2: 37, wordsBlockGroup2: 24 , ecWordsBlock: 30},
    {nversion: 29 ,correction: CorrectionH, capNum: 1677, capAlpha: 1016, capByte: 698, capKanji: 430, totalWords: 701, blocksGroup1: 19, wordsBlockGroup1: 15, blocksGroup2: 26, wordsBlockGroup2: 16 , ecWordsBlock: 30},

    {nversion: 30 ,correction: CorrectionL, capNum: 4158, capAlpha: 2520, capByte: 1732, capKanji: 1066, totalWords: 1735, blocksGroup1: 5, wordsBlockGroup1: 115, blocksGroup2: 10, wordsBlockGroup2: 116 , ecWordsBlock: 30},
    {nversion: 30 ,correction: CorrectionM, capNum: 3289, capAlpha: 1994, capByte: 1370, capKanji: 843, totalWords: 1373, blocksGroup1: 19, wordsBlockGroup1: 47, blocksGroup2: 10, wordsBlockGroup2: 48 , ecWordsBlock: 28},
    {nversion: 30 ,correction: CorrectionQ, capNum: 2358, capAlpha: 1429, capByte: 982, capKanji: 604, totalWords: 985, blocksGroup1: 15, wordsBlockGroup1: 24, blocksGroup2: 25, wordsBlockGroup2: 25 , ecWordsBlock: 30},
    {nversion: 30 ,correction: CorrectionH, capNum: 1782, capAlpha: 1080, capByte: 742, capKanji: 457, totalWords: 745, blocksGroup1: 23, wordsBlockGroup1: 15, blocksGroup2: 25, wordsBlockGroup2: 16 , ecWordsBlock: 30},

    {nversion: 31 ,correction: CorrectionL, capNum: 4417, capAlpha: 2677, capByte: 1840, capKanji: 1132, totalWords: 1843, blocksGroup1: 13, wordsBlockGroup1: 115, blocksGroup2: 3, wordsBlockGroup2: 116 , ecWordsBlock: 30},
    {nversion: 31 ,correction: CorrectionM, capNum: 3486, capAlpha: 2113, capByte: 1452, capKanji: 894, totalWords: 1455, blocksGroup1: 2, wordsBlockGroup1: 46, blocksGroup2: 29, wordsBlockGroup2: 47 , ecWordsBlock: 28},
    {nversion: 31 ,correction: CorrectionQ, capNum: 2473, capAlpha: 1499, capByte: 1030, capKanji: 634, totalWords: 1033, blocksGroup1: 42, wordsBlockGroup1: 24, blocksGroup2: 1, wordsBlockGroup2: 25 , ecWordsBlock: 30},
    {nversion: 31 ,correction: CorrectionH, capNum: 1897, capAlpha: 1150, capByte: 790, capKanji: 486, totalWords: 793, blocksGroup1: 23, wordsBlockGroup1: 15, blocksGroup2: 28, wordsBlockGroup2: 16 , ecWordsBlock: 30},

    {nversion: 32 ,correction: CorrectionL, capNum: 4686, capAlpha: 2840, capByte: 1952, capKanji: 1201, totalWords: 1955, blocksGroup1: 17, wordsBlockGroup1: 115, blocksGroup2: 0, wordsBlockGroup2: 0, ecWordsBlock: 30},
    {nversion: 32 ,correction: CorrectionM, capNum: 3693, capAlpha: 2238, capByte: 1538, capKanji: 947, totalWords: 1541, blocksGroup1: 10, wordsBlockGroup1: 46, blocksGroup2: 23, wordsBlockGroup2: 47 , ecWordsBlock: 28},
    {nversion: 32 ,correction: CorrectionQ, capNum: 2670, capAlpha: 1618, capByte: 1112, capKanji: 684, totalWords: 1115, blocksGroup1: 10, wordsBlockGroup1: 24, blocksGroup2: 35, wordsBlockGroup2: 25 , ecWordsBlock: 30},
    {nversion: 32 ,correction: CorrectionH, capNum: 2022, capAlpha: 1226, capByte: 842, capKanji: 518, totalWords: 845, blocksGroup1: 19, wordsBlockGroup1: 15, blocksGroup2: 35, wordsBlockGroup2: 16 , ecWordsBlock: 30},

    {nversion: 33 ,correction: CorrectionL, capNum: 4965, capAlpha: 3009, capByte: 2068, capKanji: 1273, totalWords: 2071, blocksGroup1: 17, wordsBlockGroup1: 115, blocksGroup2: 1, wordsBlockGroup2: 116 , ecWordsBlock: 30},
    {nversion: 33 ,correction: CorrectionM, capNum: 3909, capAlpha: 2369, capByte: 1628, capKanji: 1002, totalWords: 1631, blocksGroup1: 14, wordsBlockGroup1: 46, blocksGroup2: 21, wordsBlockGroup2: 47 , ecWordsBlock: 28},
    {nversion: 33 ,correction: CorrectionQ, capNum: 2805, capAlpha: 1700, capByte: 1168, capKanji: 719, totalWords: 1171, blocksGroup1: 29, wordsBlockGroup1: 24, blocksGroup2: 19, wordsBlockGroup2: 25 , ecWordsBlock: 30},
    {nversion: 33 ,correction: CorrectionH, capNum: 2157, capAlpha: 1307, capByte: 898, capKanji: 553, totalWords: 901, blocksGroup1: 11, wordsBlockGroup1: 15, blocksGroup2: 46, wordsBlockGroup2: 16 , ecWordsBlock: 30},

    {nversion: 34 ,correction: CorrectionL, capNum: 5253, capAlpha: 3183, capByte: 2188, capKanji: 1347, totalWords: 2191, blocksGroup1: 13, wordsBlockGroup1: 115, blocksGroup2: 6, wordsBlockGroup2: 116 , ecWordsBlock: 30},
    {nversion: 34 ,correction: CorrectionM, capNum: 4134, capAlpha: 2506, capByte: 1722, capKanji: 1060, totalWords: 1725, blocksGroup1: 14, wordsBlockGroup1: 46, blocksGroup2: 23, wordsBlockGroup2: 47 , ecWordsBlock: 28},
    {nversion: 34 ,correction: CorrectionQ, capNum: 2949, capAlpha: 1787, capByte: 1228, capKanji: 756, totalWords: 1231, blocksGroup1: 44, wordsBlockGroup1: 24, blocksGroup2: 7, wordsBlockGroup2: 25 , ecWordsBlock: 30},
    {nversion: 34 ,correction: CorrectionH, capNum: 2301, capAlpha: 1394, capByte: 958, capKanji: 590, totalWords: 961, blocksGroup1: 59, wordsBlockGroup1: 16, blocksGroup2: 1, wordsBlockGroup2: 17 , ecWordsBlock: 30},

    {nversion: 35 ,correction: CorrectionL, capNum: 5529, capAlpha: 3351, capByte: 2303, capKanji: 1417, totalWords: 2306, blocksGroup1: 12, wordsBlockGroup1: 121, blocksGroup2: 7, wordsBlockGroup2: 122 , ecWordsBlock: 30},
    {nversion: 35 ,correction: CorrectionM, capNum: 4343, capAlpha: 2632, capByte: 1809, capKanji: 1113, totalWords: 1812, blocksGroup1: 12, wordsBlockGroup1: 47, blocksGroup2: 26, wordsBlockGroup2: 48 , ecWordsBlock: 28},
    {nversion: 35 ,correction: CorrectionQ, capNum: 3081, capAlpha: 1867, capByte: 1283, capKanji: 790, totalWords: 1286, blocksGroup1: 39, wordsBlockGroup1: 24, blocksGroup2: 14, wordsBlockGroup2: 25 , ecWordsBlock: 30},
    {nversion: 35 ,correction: CorrectionH, capNum: 2361, capAlpha: 1431, capByte: 983, capKanji: 605, totalWords: 986, blocksGroup1: 22, wordsBlockGroup1: 15, blocksGroup2: 41, wordsBlockGroup2: 16 , ecWordsBlock: 30},

    {nversion: 36 ,correction: CorrectionL, capNum: 5836, capAlpha: 3537, capByte: 2431, capKanji: 1496, totalWords: 2434, blocksGroup1: 6, wordsBlockGroup1: 121, blocksGroup2: 14, wordsBlockGroup2: 122 , ecWordsBlock: 30},
    {nversion: 36 ,correction: CorrectionM, capNum: 4588, capAlpha: 2780, capByte: 1911, capKanji: 1176, totalWords: 1914, blocksGroup1: 6, wordsBlockGroup1: 47, blocksGroup2: 34, wordsBlockGroup2: 48 , ecWordsBlock: 28},
    {nversion: 36 ,correction: CorrectionQ, capNum: 3244, capAlpha: 1966, capByte: 1351, capKanji: 832, totalWords: 1354, blocksGroup1: 46, wordsBlockGroup1: 24, blocksGroup2: 10, wordsBlockGroup2: 25 , ecWordsBlock: 30},
    {nversion: 36 ,correction: CorrectionH, capNum: 2524, capAlpha: 1530, capByte: 1051, capKanji: 647, totalWords: 1054, blocksGroup1: 2, wordsBlockGroup1: 15, blocksGroup2: 64, wordsBlockGroup2: 16 , ecWordsBlock: 30},

    {nversion: 37 ,correction: CorrectionL, capNum: 6153, capAlpha: 3729, capByte: 2563, capKanji: 1577, totalWords: 2566, blocksGroup1: 17, wordsBlockGroup1: 122, blocksGroup2: 4, wordsBlockGroup2: 123 , ecWordsBlock: 30},
    {nversion: 37 ,correction: CorrectionM, capNum: 4775, capAlpha: 2894, capByte: 1989, capKanji: 1224, totalWords: 1992, blocksGroup1: 29, wordsBlockGroup1: 46, blocksGroup2: 14, wordsBlockGroup2: 47 , ecWordsBlock: 28},
    {nversion: 37 ,correction: CorrectionQ, capNum: 3417, capAlpha: 2071, capByte: 1423, capKanji: 876, totalWords: 1426, blocksGroup1: 49, wordsBlockGroup1: 24, blocksGroup2: 10, wordsBlockGroup2: 25 , ecWordsBlock: 30},
    {nversion: 37 ,correction: CorrectionH, capNum: 2625, capAlpha: 1591, capByte: 1093, capKanji: 673, totalWords: 1096, blocksGroup1: 24, wordsBlockGroup1: 15, blocksGroup2: 46, wordsBlockGroup2: 16 , ecWordsBlock: 30},

    {nversion: 38 ,correction: CorrectionL, capNum: 6479, capAlpha: 3927, capByte: 2699, capKanji: 1661, totalWords: 2702, blocksGroup1: 4, wordsBlockGroup1: 122, blocksGroup2: 18, wordsBlockGroup2: 123 , ecWordsBlock: 30},
    {nversion: 38 ,correction: CorrectionM, capNum: 5039, capAlpha: 3054, capByte: 2099, capKanji: 1292, totalWords: 2102, blocksGroup1: 13, wordsBlockGroup1: 46, blocksGroup2: 32, wordsBlockGroup2: 47 , ecWordsBlock: 28},
    {nversion: 38 ,correction: CorrectionQ, capNum: 3599, capAlpha: 2181, capByte: 1499, capKanji: 923, totalWords: 1502, blocksGroup1: 48, wordsBlockGroup1: 24, blocksGroup2: 14, wordsBlockGroup2: 25 , ecWordsBlock: 30},
    {nversion: 38 ,correction: CorrectionH, capNum: 2735, capAlpha: 1658, capByte: 1139, capKanji: 701, totalWords: 1142, blocksGroup1: 42, wordsBlockGroup1: 15, blocksGroup2: 32, wordsBlockGroup2: 16 , ecWordsBlock: 30},

    {nversion: 39 ,correction: CorrectionL, capNum: 6743, capAlpha: 4087, capByte: 2809, capKanji: 1729, totalWords: 2812, blocksGroup1: 20, wordsBlockGroup1: 117, blocksGroup2: 4, wordsBlockGroup2: 118 , ecWordsBlock: 30},
    {nversion: 39 ,correction: CorrectionM, capNum: 5313, capAlpha: 3220, capByte: 2213, capKanji: 1362, totalWords: 2216, blocksGroup1: 40, wordsBlockGroup1: 47, blocksGroup2: 7, wordsBlockGroup2: 48 , ecWordsBlock: 28},
    {nversion: 39 ,correction: CorrectionQ, capNum: 3791, capAlpha: 2298, capByte: 1579, capKanji: 972, totalWords: 1582, blocksGroup1: 43, wordsBlockGroup1: 24, blocksGroup2: 22, wordsBlockGroup2: 25 , ecWordsBlock: 30},
    {nversion: 39 ,correction: CorrectionH, capNum: 2927, capAlpha: 1774, capByte: 1219, capKanji: 750, totalWords: 1222, blocksGroup1: 10, wordsBlockGroup1: 15, blocksGroup2: 67, wordsBlockGroup2: 16 , ecWordsBlock: 30},

    {nversion: 40 ,correction: CorrectionL, capNum: 7089, capAlpha: 4296, capByte: 2953, capKanji: 1817, totalWords: 2956, blocksGroup1: 19, wordsBlockGroup1: 118, blocksGroup2: 6, wordsBlockGroup2: 119 , ecWordsBlock: 30},
    {nversion: 40 ,correction: CorrectionM, capNum: 5596, capAlpha: 3391, capByte: 2331, capKanji: 1435, totalWords: 2334, blocksGroup1: 18, wordsBlockGroup1: 47, blocksGroup2: 31, wordsBlockGroup2: 48 , ecWordsBlock: 28},
    {nversion: 40 ,correction: CorrectionQ, capNum: 3993, capAlpha: 2420, capByte: 1663, capKanji: 1024, totalWords: 1666, blocksGroup1: 34, wordsBlockGroup1: 24, blocksGroup2: 34, wordsBlockGroup2: 25 , ecWordsBlock: 30},
    {nversion: 40 ,correction: CorrectionH, capNum: 3057, capAlpha: 1852, capByte: 1273, capKanji: 784, totalWords: 1276, blocksGroup1: 20, wordsBlockGroup1: 15, blocksGroup2: 61, wordsBlockGroup2: 16 , ecWordsBlock: 30},
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


func errorCorrection(encoded []byte, version Version) {
  fmt.Printf("coefficients: %d\n", encoded)
}


