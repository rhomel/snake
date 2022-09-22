package main

import (
	"image/color"
	"log"
	"math/rand"
	"os"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/rhomel/snake/pkg/images"
)

const (
	title = "Snake"

	boardSize = 20

	// screen dimensions
	screenWidth  = boardSize * (images.TileWidth + 2)
	screenHeight = boardSize * (images.TileHeight + 2)

	// window dimensions
	windowWidth  = screenWidth * 2
	windowHeight = screenHeight * 2
)

var (
	// colors
	lightGreen       = color.RGBA{0, 0xc1, 0x62, 0xff}
	lightPastelGreen = color.RGBA{0x79, 0xd2, 0xa7, 0xff}
	darkGreen        = color.RGBA{0, 0x97, 0x60, 0xff}
	darkGreen2       = color.RGBA{0x26, 0x86, 0x63, 0xff}
	green            = color.RGBA{0x00, 0x8b, 0x02, 0xff}
	brown            = color.RGBA{0x8b, 0x57, 0x2a, 0xff}

	// images
	whiteTile = images.NewTile()
	emptyTile = images.NewColoredTile(lightGreen)
	foodTile  = images.NewColoredTile(brown)
	headTile  = images.NewColoredTile(green)
	bodyTile  = images.NewColoredTile(darkGreen)
)

type cell string

const (
	empty cell = "empty"
	food  cell = "food"
	head  cell = "head"
	tail  cell = "tail"
	body  cell = "body"
)

func (c cell) Tile() *ebiten.Image {
	switch c {
	case empty:
		return emptyTile
	case food:
		return foodTile
	default:
		return whiteTile
	}
}

type board struct {
	cells [boardSize][boardSize]cell
	food  int
}

func emptyBoard() *board {
	b := &board{
		cells: [boardSize][boardSize]cell{},
	}
	for x := 0; x < boardSize; x++ {
		for y := 0; y < boardSize; y++ {
			b.cells[x][y] = empty
		}
	}
	return b
}

func translateByTileSize(g ebiten.GeoM, x, y int) ebiten.GeoM {
	g.Translate(float64(x)*images.TileWidth, float64(y)*images.TileHeight)
	return g
}

func (b *board) Draw(screen *ebiten.Image, origin ebiten.GeoM) {
	for x := 0; x < boardSize; x++ {
		for y := 0; y < boardSize; y++ {
			screen.DrawImage(b.cells[x][y].Tile(), &ebiten.DrawImageOptions{
				GeoM: translateByTileSize(origin, x, y),
			})
		}
	}
}

func (b *board) AddFood(player *snake) {
	if b.food >= 5 {
		return
	}
	var x, y int
	for tries := 10; tries > 0; tries-- {
		x, y = b.randomPoint()
		if !player.IsAt(x, y) {
			break
		}
	}
	b.cells[x][y] = food
	b.food++
}

func (b *board) RemoveFood(x, y int) {
	if b.IsFood(x, y) {
		b.cells[x][y] = empty
		b.food--
	}
}

func (b *board) IsValidPosition(x, y int) bool {
	return x >= 0 && x < boardSize && y >= 0 && y < boardSize
}

func (b *board) IsFood(x, y int) bool {
	if !b.IsValidPosition(x, y) {
		return false
	}
	return b.cells[x][y] == food
}

func (b *board) randomPoint() (int, int) {
	x := rand.Intn(boardSize)
	y := rand.Intn(boardSize)
	return x, y
}

type position struct {
	x int
	y int
}

func (p position) Draw(screen, tile *ebiten.Image, origin ebiten.GeoM) {
	screen.DrawImage(tile, &ebiten.DrawImageOptions{
		GeoM: translateByTileSize(origin, p.x, p.y),
	})
}

func (p position) IsOnBoard() bool {
	return p.x >= 0 && p.x < boardSize && p.y >= 0 && p.y < boardSize
}

type direction int

const (
	left  direction = 1
	right direction = 2
	up    direction = 3
	down  direction = 4
)

func (d direction) String() string {
	switch d {
	case left:
		return "left"
	case right:
		return "right"
	case up:
		return "up"
	case down:
		return "down"
	default:
		return "unknown direction"
	}
}

const capacity = boardSize * boardSize

type ring struct {
	positions [capacity]position

	size int
	h    int // head index
	t    int // tail index
}

func newRing(head position) *ring {
	r := &ring{
		positions: [capacity]position{head},
		h:         0,
		t:         0,
		size:      1,
	}
	return r
}

func (s *ring) head() position {
	return s.positions[s.h]
}

func (s *ring) tail() position {
	return s.positions[s.t]
}

// move makes position p the new head and shifts all segments by one removing
// the last position. The last position is returned.
func (s *ring) move(p position) position {
}

func (s *ring) grow(p position) {
}

type snake struct {
	direction direction

	last position
	body []position

	alive bool
}

func NewSnake() *snake {
	center := int(boardSize / 2)
	head := position{center, center}
	s := &snake{
		direction: left,
		body:      []position{head},
		last:      head,
		alive:     true,
	}
	return s
}

func (s *snake) Size() int {
	return len(s.body)
}

func (s *snake) IsAt(x, y int) bool {
	head := s.head()
	if head.x == x && head.y == y {
		return true
	}
	if s.IsBodyAt(x, y) {
		return true
	}
	return false
}

func (s *snake) IsBodyAt(x, y int) bool {
	for i := 1; i < len(s.body); i++ {
		if s.body[i].x == x && s.body[i].y == y {
			return true
		}
	}
	return false
}

func (s *snake) Draw(screen *ebiten.Image, origin ebiten.GeoM) {
	for _, bodySegment := range s.body {
		bodySegment.Draw(screen, bodyTile, origin)
	}
	s.head().Draw(screen, headTile, origin)
}

func (s *snake) SetDirection(d direction) {
	if len(s.body) > 1 && isOppositeDirection(s.direction, d) {
		// reject going the direction of the body if longer than 1
		return
	}
	s.direction = d
}

func isOppositeDirection(current, desired direction) bool {
	if current == left && desired == right {
		return true
	}
	if current == right && desired == left {
		return true
	}
	if current == up && desired == down {
		return true
	}
	if current == down && desired == up {
		return true
	}
	return false
}

// slither advances the snake one position in the current direction
func (s *snake) slither() bool {
	next := s.getNextPosition()
	if !next.IsOnBoard() {
		return false
	}
	s.last = s.tail()
	s.body = append([]position{next}, s.body[:len(s.body)-1]...)
	if s.IsBodyAt(next.x, next.y) {
		return false
	}
	return true
}

func (s *snake) move(board *board) {
	s.alive = s.slither()
	head := s.head()
	if board.IsFood(head.x, head.y) {
		board.RemoveFood(head.x, head.y)
		s.grow()
	}
}

func (s *snake) grow() {
	s.body = append(s.body, s.last)
}

func (s *snake) head() position {
	return s.body[0]
}

func (s *snake) tail() position {
	return s.body[len(s.body)-1]
}

func (s *snake) getNextPosition() position {
	next := s.head()
	switch s.direction {
	case left:
		next.x--
	case right:
		next.x++
	case up:
		next.y--
	case down:
		next.y++
	}
	return next
}

type state int

const (
	gameover state = 0
	playing  state = 1
)

type Game struct {
	state  state
	board  *board
	player *snake
	ticks  int
}

func NewGame() *Game {
	return &Game{
		state:  playing,
		board:  emptyBoard(),
		player: NewSnake(),
	}
}

func (game *Game) Restart() {
	game.board = emptyBoard()
	game.player = NewSnake()
	game.state = playing
}

func (game *Game) DrawScore(screen *ebiten.Image) {
	score := strconv.Itoa(game.player.Size())
	padding := 5
	x := padding
	y := screenHeight - images.TileHeight - padding
	ebitenutil.DebugPrintAt(screen, score, x, y)
}

func (game *Game) Draw(screen *ebiten.Image) {
	screen.Fill(lightPastelGreen)
	boardOrigin := ebiten.GeoM{}
	boardOrigin.Translate(20, 20)
	game.board.Draw(screen, boardOrigin)
	game.player.Draw(screen, boardOrigin)
	game.DrawScore(screen)
	if game.state == gameover {
		ebitenutil.DebugPrint(screen, "press space to restart")
	}
}

func (*Game) Layout(outsideWidth int, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (game *Game) UpdateBoard() {
	if game.ticks%10 == 0 {
		game.board.AddFood(game.player)
	}
}

func (game *Game) UpdatePlay() {
	check := inpututil.IsKeyJustPressed
	switch {
	case check(ebiten.KeyArrowUp):
		game.player.SetDirection(up)
	case check(ebiten.KeyW):
		game.player.SetDirection(up)
	case check(ebiten.KeyArrowLeft):
		game.player.SetDirection(left)
	case check(ebiten.KeyA):
		game.player.SetDirection(left)
	case check(ebiten.KeyArrowDown):
		game.player.SetDirection(down)
	case check(ebiten.KeyS):
		game.player.SetDirection(down)
	case check(ebiten.KeyArrowRight):
		game.player.SetDirection(right)
	case check(ebiten.KeyD):
		game.player.SetDirection(right)
	}
	game.player.move(game.board)
	if !game.player.alive {
		game.state = gameover
	}
	game.ticks++
	game.UpdateBoard()
}

func (game *Game) Update() error {
	check := inpututil.IsKeyJustPressed
	if game.state == playing {
		game.UpdatePlay()
	}
	if game.state == gameover {
		if check(ebiten.KeySpace) {
			game.Restart()
		}
	}
	if check(ebiten.KeyEscape) {
		log.Println("quit")
		os.Exit(0)
	}
	//log.Println("direction:", game.player.direction.String())
	return nil
}

func main() {
	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowTitle(title)
	ebiten.SetTPS(10)
	game := NewGame()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
