# XVEC — Vector Graphics Format

## 1. Introduction

### 1.1. Purpose

XVEC is a plain‑text vector graphics format designed for low‑resolution
and pixel‑art. It stores a fixed‑size canvas, a set of drawing instructions,
and an anti‑aliasing flag.

### 1.2. Terminology

The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT", "SHOULD",
"SHOULD NOT", "RECOMMENDED", "MAY", and "OPTIONAL" in this document are to be
interpreted as described in [RFC 2119](https://www.rfc-editor.org/rfc/rfc2119).

## 2. File Structure

### 2.1. Character Encoding and Tokenisation

The file is UTF‑8 text.  A parser MUST tokenise the input into a sequence of
tokens, where each token is one of:

- A **keyword** — an alphabetic identifier (`xvec`, `size`, `circle`, `disk`,
  `rect`, `slab`, `line`, `fill`, `stroke`, `end`, `move`, `quad`, `cubic`,
  `arc`, `arcto`, `close`, `antialias`, `true`, `false`, `C`, `CC`,
  `join`, `miter`, `bevel`, `round`, `cap`, `butt`, `square`, `rule`,
  `evenodd`, `nonzero`).
- A **number** — an integer (`42`) or floating‑point (`3.14`, `.5`, `1e2`)
  literal.
- A **colour** — a `#` character followed by exactly eight hexadecimal digits
  (`#ff0000ff`).
- **Punctuation** — the keyword `line` is used both as a top‑level instruction
  and as a path step inside `fill`/`stroke` blocks (see Section 4.5).

Whitespace (spaces, tabs, newlines, carriage returns) separates tokens and is
otherwise ignored.  Commas are also treated as whitespace.

Both `//` line comments (from `//` to end of line) and `/* */` block comments
MUST be recognised and skipped.  Comments MAY appear anywhere a token is
expected, including between instructions, after values on the same line, and
before the `xvec` header.

### 2.2. Grammar Notation

```
CAPITALISED  = terminal keyword
<lowercase>  = non‑terminal
`"literal"'  = literal character
( )          = grouping
[ ]          = optional (zero or one)
{ }          = repetition (zero or more)
|            = alternation
```

### 2.3. Overall Structure

```
xvec 1
size <width> <height>
[antialias <bool>]
{ <instruction> }
```

The first non‑comment token MUST be the literal keyword `xvec` followed by the
version token `1`.  Any other version string is an error.

### 2.4. The `<number>` Token

A `<number>` is a token that represents an integer or floating‑point value
parsable as a 32‑bit IEEE 754 floating‑point number.  Both integer (`42`) and
floating‑point (`3.14`, `.5`, `1e2`) forms are accepted.  Parsers MUST reject
tokens that cannot be parsed as `float32`.

### 2.5. The `<colour>` Token

A `<colour>` is exactly eight hexadecimal digits (`0–9`, `a–f`, `A–F`) prefixed
with `#`, representing the colour in big‑endian **RRGGBBAA** order.

Examples:

```
#ff0000ff   -- opaque red
#00ff0080   -- semi‑transparent green
#00000000   -- fully transparent black
```

Implementations MUST reject a `<colour>` whose first token is not `#`.
Implementations MUST reject hex strings longer or shorter than eight digits.

### 2.6. The `<bool>` Token

A `<bool>` is the identifier `true` or `false`.  Any other value is an error.

## 3. Global Declarations

### 3.1. Header

```
xvec 1
```

REQUIRED.  MUST appear exactly once as the first meaningful token of the file.
MUST be followed by the version number `1` as a single separate token.  Future
versions of the format will increment this number; a parser MUST reject an
unrecognised version.

### 3.2. Canvas Size

```
size <width> <height>
```

REQUIRED.  MUST appear exactly once.  `<width>` and `<height>` are `<number>`
tokens giving the logical canvas dimensions in pixels.  If absent, a parser
SHOULD default to `320 240`.

### 3.3. Anti‑Aliasing

```
antialias <bool>
```

OPTIONAL.  Controls whether path and primitive rendering uses anti‑aliasing.
If absent, a parser SHOULD default to `true`.

## 4. Drawing Instructions

### 4.1. Circle (Stroke)

```
circle <cx> <cy> <r> <stroke> <colour>
```

Draws the outline of a circle centred at (`<cx>`, `<cy>`) with radius `<r>`
and stroke width `<stroke>`.  All values are `<number>` tokens; `<colour>` is
a `<colour>` token.

### 4.2. Disk (Fill)

```
disk <cx> <cy> <r> <colour>
```

Draws a filled circle centred at (`<cx>`, `<cy>`) with radius `<r>`.  All
values are `<number>` tokens; `<colour>` is a `<colour>` token.

### 4.3. Rectangle (Stroke)

```
rect <x> <y> <w> <h> <stroke> <colour>
```

Draws the outline of a rectangle with top‑left corner (`<x>`, `<y>`), width
`<w>`, height `<h>`, and stroke width `<stroke>`.

### 4.4. Filled Rectangle (Slab)

```
slab <x> <y> <w> <h> <colour>
```

Draws a filled rectangle with top‑left corner (`<x>`, `<y>`), width `<w>`,
height `<h>`.

### 4.5. Line Segment

```
line <x1> <y1> <x2> <y2> <stroke> <colour>
```

Draws a straight line segment from (`<x1>`, `<y1>`) to (`<x2>`, `<y2>`) with
stroke width `<stroke>`.

When the token `line` appears *inside* a `fill`/`stroke` path block (see
Section 5), it is parsed as a path step rather than a standalone instruction.

### 4.6. Filled Path

```
fill <colour>
  { <path‑step> }
end
```

Fills a closed vector path built from the enclosed `<path‑step>` lines.
If the path block contains at least one step, the final step MUST be `close`.
An empty path block (no steps between `fill` and `end`) is allowed.

### 4.7. Stroked Path

```
stroke <width> <colour>
  { <path‑step> }
end
```

Strokes a vector path built from the enclosed `<path‑step>` lines with stroke
width `<width>`.  If the path block contains at least one step, the final step
MUST be `close`.  An empty path block is allowed.

## 5. Path Steps And Options

The following steps and options are recognised only inside `fill` ... `end` or
`stroke` ... `end` blocks.

### 5.1. MoveTo

```
move <x> <y>
```

Lifts the pen and moves to (`<x>`, `<y>`).  Begins a new sub‑path.

### 5.2. LineTo

```
line <x> <y>
```

Draws a straight line from the current pen position to (`<x>`, `<y>`).

### 5.3. QuadTo (Quadratic Bézier)

```
quad <x1> <y1> <x2> <y2>
```

Draws a quadratic Bézier curve from the current pen position to (`<x2>`, `<y2>`)
using (`<x1>`, `<y1>`) as the control point.

### 5.4. CubicTo (Cubic Bézier)

```
cubic <x1> <y1> <x2> <y2> <x3> <y3>
```

Draws a cubic Bézier curve from the current pen position to (`<x3>`, `<y3>`)
using (`<x1>`, `<y1>`) and (`<x2>`, `<y2>`) as control points.

### 5.5. Arc

```
arc <cx> <cy> <r> <start> <end> <dir>
```

Draws a circular arc centred at (`<cx>`, `<cy>`) with radius `<r>`, starting
at angle `<start>` (radians) and ending at `<end>` (radians).

`<dir>` is the identifier `C` (clockwise) or `CC` (counter‑clockwise).

### 5.6. ArcTo

```
arcto <x1> <y1> <x2> <y2> <r>
```

Draws a circular arc from the current pen position to (`<x2>`, `<y2>`), with
a turning point at (`<x1>`, `<y1>`) and radius `<r>`.

### 5.7. Close

```
close
```

Closes the current sub‑path by drawing a straight line back to the most recent
`move` point.


### 5.8. Rule

```
rule [evenodd|nonzero]
```

Rule option sets the fill rule of the fill block it is in to even/odd filling
or nonzero filling. It may not be used in a stroke block.

### 5.9. Cap

```
cap [butt|round|square]
```

Cap option sets the line caps of the stroke block it is in to
butt, round, or square. It may not be used in a stroke block.

### 5.6. Join


```
join [miter|round|bevel]
```

Join option sets the line joints of the stroke block it is in to
miter, round, or bevel. It may not be used in a stroke block.


## 6. Sub‑Paths

A `fill` or `stroke` path MAY contain multiple sub‑paths.  Each sub‑path is
begun by a `move` step and MAY be ended by a `close` step.  The overall step
sequence MUST end with `close` if non‑empty.

```
fill #ff0000ff
  move 10 10
  line 90 10
  line 50 40
  close
  move 10 60
  line 90 60
  line 50 90
  close
end
```

The example above fills two separate triangular regions in a single draw
operation.

## 7. Comments

Two comment forms are supported:

- `//` — line comment; all text from `//` to the end of the line is ignored.
- `/* ... */` — block comment; all text between `/*` and `*/` is ignored.
  Block comments MAY span multiple lines and MUST NOT nest.

Comments MAY appear anywhere a token is expected, including between
instructions, after values on the same line, and before the `xvec` header.

## 8. Full Example

```
xvec 1
size 160 120
antialias true

slab 0 0 160 120 #1e1e32ff

// Outlined circle
circle 80 60 50 2 #ffffffff

// Filled circle
disk 80 60 20 #ff0000ff

// Stroked rectangle
rect 10 10 60 40 1 #00ff00ff

// Line segment
line 0 0 160 120 1 #0000ffff

// Filled path with two sub‑paths
fill #00c8c864
  move 80 20
  line 140 60
  line 80 100
  close
  move 20 60
  line 60 60
  line 40 80
  close
end

/* A stroked cubic Bézier */
stroke 1 #ffff00ff
  move 40 30
  cubic 120 10 120 110 40 90
  close
end
```

## 9. Error Handling

A parser MUST abort and report an error when it encounters:

- An unsupported `xvec` version.
- A `<colour>` token that does not begin with `#`.
- A `<number>` token that cannot be parsed as `float32`.
- An unrecognised keyword at the top level or inside a path block.
- A `fill` or `stroke` block with one or more steps whose last step is not
  `close`.
- A path step keyword (`move`, `line`, `quad`, `cubic`, `arc`, `arcto`,
  `close`) appearing outside a `fill`/`stroke` block.
- Unexpected end of input while reading a required value.
