/**********************************************************\
|                                                          |
|                          hprose                          |
|                                                          |
| Official WebSite: http://www.hprose.com/                 |
|                   http://www.hprose.org/                 |
|                                                          |
\**********************************************************/
/**********************************************************\
 *                                                        *
 * io/tags.go                                             *
 *                                                        *
 * hprose tags enum for Go.                               *
 *                                                        *
 * LastModified: Aug 15, 2016                             *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package io

// Hprose Tags
const (

	// Serialize Type
	TagInteger  byte = 'i'
	TagLong     byte = 'l'
	TagDouble   byte = 'd'
	TagNull     byte = 'n'
	TagEmpty    byte = 'e'
	TagTrue     byte = 't'
	TagFalse    byte = 'f'
	TagNaN      byte = 'N'
	TagInfinity byte = 'I'
	TagDate     byte = 'D'
	TagTime     byte = 'T'
	TagUTC      byte = 'Z'
	TagBytes    byte = 'b'
	TagUTF8Char byte = 'u'
	TagString   byte = 's'
	TagGUID     byte = 'g'
	TagList     byte = 'a'
	TagMap      byte = 'm'
	TagClass    byte = 'c'
	TagObject   byte = 'o'
	TagRef      byte = 'r'

	// Serialize Marks
	TagPos        byte = '+'
	TagNeg        byte = '-'
	TagSemicolon  byte = ';'
	TagOpenbrace  byte = '{'
	TagClosebrace byte = '}'
	TagQuote      byte = '"'
	TagPoint      byte = '.'

	// Protocol Tags
	TagFunctions byte = 'F'
	TagCall      byte = 'C'
	TagResult    byte = 'R'
	TagArgument  byte = 'A'
	TagError     byte = 'E'
	TagEnd       byte = 'z'
)
