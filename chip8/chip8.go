package chip8

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"time"
)

var fontset = [...]byte{
	0xF0, 0x90, 0x90, 0x90, 0xF0, //0
	0x20, 0x60, 0x20, 0x20, 0x70, //1
	0xF0, 0x10, 0xF0, 0x80, 0xF0, //2
	0xF0, 0x10, 0xF0, 0x10, 0xF0, //3
	0x90, 0x90, 0xF0, 0x10, 0x10, //4
	0xF0, 0x80, 0xF0, 0x10, 0xF0, //5
	0xF0, 0x80, 0xF0, 0x90, 0xF0, //6
	0xF0, 0x10, 0x20, 0x40, 0x40, //7
	0xF0, 0x90, 0xF0, 0x90, 0xF0, //8
	0xF0, 0x90, 0xF0, 0x10, 0xF0, //9
	0xF0, 0x90, 0xF0, 0x90, 0x90, //A
	0xE0, 0x90, 0xE0, 0x90, 0xE0, //B
	0xF0, 0x80, 0x80, 0x80, 0xF0, //C
	0xE0, 0x90, 0x90, 0x90, 0xE0, //D
	0xF0, 0x80, 0xF0, 0x80, 0xF0, //E
	0xF0, 0x80, 0xF0, 0x80, 0x80, //F
}

type Chip8 struct {

	DrawFlag bool
	Running bool

	opcode uint16		// Current Opcode
	memory [4096]byte	// Memory (size = 4k)
	registers [16]byte	// V-regs (V0 - VF)
	index int			// Index Register
	pc int				// Program Counter
	gfx [64 * 32]byte	// Graphics Memory (2048 pixels)
	
	stack [16]int		// Stack (16 levels)
	sp int				// Stack Pointer
	
	delayTimer byte
	soundTimer byte

	key [16]byte		// Input Buttons
}

func NewChip8() (*Chip8) {
	cpu := new(Chip8)
	cpu.Init()
	return cpu
}

func (chip *Chip8) Init() {
	chip.pc = 0x200
	chip.opcode = 0
	chip.index = 0
	chip.Running = true

	// Clear memory
	for i := 0; i < 4096; i++ {
		chip.memory[i] = 0
	}

	// Load fontset
	copy(chip.memory[0:80], fontset[:])

	rand.Seed(time.Now().Unix())
}

func (chip *Chip8) LoadRom(filePath string) {

	fmt.Printf("Loading file: '%s'\n", filePath)

	bytes, err := ioutil.ReadFile(filePath)

	if err != nil {
		panic(err)
	}

	// Rom cartridge starts at location 0x0100
	copy(chip.memory[0x200:], bytes[:])

}

func (chip *Chip8) Dump() {
	for i := 0; i < len(chip.memory); i++ {

		hexChar := "%02x"
		if i % 8 == 0 {
			fmt.Printf("%04x\t", i)
			hexChar += " "
		} else if (i + 1) % 8 == 0 {
			hexChar += "\n"
		} else {
			hexChar += " "
		}

		fmt.Printf(hexChar, chip.memory[i])
	}
}

func (chip *Chip8) EmulateCycle() {

	// Fetch opcode
	chip.opcode = uint16((int(chip.memory[chip.pc]) << 8) | int(chip.memory[chip.pc + 1]))
	
	// Execute opcode
	opcodes[chip.memory[chip.pc] & 0xF0](chip, chip.opcode)

	// Update timers
	if chip.delayTimer > 0 {
		chip.delayTimer--
	}

	if chip.soundTimer > 0 {
		if chip.soundTimer == 1 {
			fmt.Println("BEEP!")
		}

		chip.soundTimer--
	}

}

func (chip *Chip8) DumpScreen() {

	screen := ""

	for i := 0; i < 32; i++ {
		for j := 0; j < 64; j++ {
			if(chip.gfx[i*64 +j] == 1) {
				screen += "X"
			} else {
				screen += " "
			}
		}
		screen += "\n"
	}

	fmt.Printf(screen)
	chip.DrawFlag = false
}

func (chip *Chip8) ToString() string {
	
	return fmt.Sprintf("Chip8 {\n\tpc:\t%04x,\n\tsp:\t%v,\n\tStack:\t%0 4x,\n\tV:\t%0 2x,\n\tI:\t%02x,\n}", chip.pc, chip.sp, chip.stack, chip.registers, chip.index)

}

// -- Opcode map --
// ex. opcodes[0x00](chip);
var opcodes = map[byte]func(*Chip8, uint16){
	0x00: hex00,
	0x10: hex10,
	0x20: hex20,
	0x30: hex30,
	0x40: hex40,
	0x60: hex60,
	0x70: hex70,
	0x80: hex80,
	0xA0: hexA0,
	0xB0: hexB0,
	0xC0: hexC0,
	0xD0: hexD0,
	0xE0: hexE0,
	0xF0: hexF0,
}

func hex00(chip *Chip8, opcode uint16) {
	switch opcode & 0x000F {
		case 0x0000: // Clears the screen
			for i := 0; i < 2048; i++ {
				chip.gfx[i] = 0
			}
			chip.DrawFlag = true
			chip.pc += 2

		case 0x000E: // Returns from subroutine
			chip.sp--
			chip.pc = int(chip.stack[chip.sp])
			chip.pc += 2
	}
}

/**
 * 1NNN - Jumps to address NNN
 */
func hex10(chip *Chip8, opcode uint16) {
	chip.pc = int(opcode & 0x0FFF)
}

/**
 * 2NNN - Calls subroutine at NNN
 */
func hex20(chip *Chip8, opcode uint16) {
	chip.stack[chip.sp] = chip.pc	// Store current address in stack
	chip.sp++							// Increment stack pointer
	chip.pc = int(opcode & 0x0FFF)		// Set the program counter to the address at NNN
}

/**
 * 3XNN - Skips the next instruction if VX equals NN
 */
func hex30(chip *Chip8, opcode uint16) {
	if chip.registers[(opcode & 0x0F00) >> 8] == byte(opcode & 0x00FF) {
		chip.pc += 4
	} else {
		chip.pc += 2
	}
}

/**
 * 4XNN - Skips the next instruction if VX doesn't equal NN
 */
func hex40(chip *Chip8, opcode uint16) {
	if chip.registers[(opcode & 0x0F00) >> 8] != byte(opcode & 0x00FF) {
		chip.pc += 4
	} else {
		chip.pc += 2
	}
}

/**
 * 6XNN - Sets register X to NN
 */
func hex60(chip *Chip8, opcode uint16) {
	chip.registers[(opcode & 0x0F00) >> 8] = byte(opcode & 0x00FF)
	chip.pc += 2
}

/**
 * 7XNN	- Adds NN to VX.
 */
func hex70(chip *Chip8, opcode uint16) {
	chip.registers[(opcode & 0x0F00) >> 8] += byte(opcode & 0x00FF)
	chip.pc += 2
}

/**
 * 8XY?
 */
func hex80(chip *Chip8, opcode uint16) {
	switch opcode & 0x000F {
		case 0x0000: // 0x8XY0: Sets VX to the value of VY
			chip.registers[(opcode & 0x0F00) >> 8] = chip.registers[(opcode & 0x00F0) >> 4]
			chip.pc += 2
		case 0x0001: // 0x8XY1: Sets VX to "VX OR VY"
			chip.registers[(opcode & 0x0F00) >> 8] |= chip.registers[(opcode & 0x00F0) >> 4] 
			chip.pc += 2
		case 0x0002: // 0x8XY2: Sets VX to "VX AND VY"
			chip.registers[(opcode & 0x0F00) >> 8] &= chip.registers[(opcode & 0x00F0) >> 4] 
			chip.pc += 2
		case 0x0003: // 0x8XY3: Sets VX to "VX XOR VY"
			chip.registers[(opcode & 0x0F00) >> 8] ^= chip.registers[(opcode & 0x00F0) >> 4] 
			chip.pc += 2
		case 0x0004: // 0x8XY4: Adds VY to VX. VF is set to 1 when there's a carry, and to 0 when there isn't
			if chip.registers[(opcode & 0x00F0) >> 4]  > (0xFF - chip.registers[(opcode & 0x0F00) >> 8]) {
				chip.registers[0xF] = 1 // carry
			} else {
				chip.registers[0xF] = 0
			}
			chip.registers[(opcode & 0x0F00) >> 8] += chip.registers[(opcode & 0x00F0) >> 4]
			chip.pc += 2
		case 0x0005: // 0x8XY5: VY is subtracted from VX. VF is set to 0 when there's a borrow, and 1 when there isn't
			if chip.registers[(opcode & 0x00F0) >> 4]  > chip.registers[(opcode & 0x0F00) >> 8] {
				chip.registers[0xF] = 0; // there is a borrow
			} else {
				chip.registers[0xF] = 1
			}
			chip.registers[(opcode & 0x0F00) >> 8] -= chip.registers[(opcode & 0x00F0) >> 4]
			chip.pc += 2
		case 0x0006: // 0x8XY6: Shifts VX right by one. VF is set to the value of the least significant bit of VX before the shift
			chip.registers[0xF] = chip.registers[(opcode & 0x0F00) >> 8] & 0x1
			chip.registers[(opcode & 0x0F00) >> 8] >>= 1
			chip.pc += 2
		case 0x0007: // 0x8XY7: Sets VX to VY minus VX. VF is set to 0 when there's a borrow, and 1 when there isn't
			if chip.registers[(opcode & 0x0F00) >> 8] > chip.registers[(opcode & 0x00F0) >> 4] {
				chip.registers[0xF] = 0 // There is a borrow
			} else {
				chip.registers[0xF] = 1
			}
			chip.registers[(opcode & 0x0F00) >> 8] = chip.registers[(opcode & 0x00F0) >> 4] - chip.registers[(opcode & 0x0F00) >> 8]
			chip.pc += 2
		case 0x000E: // 0x8XYE: Shifts VX left by one. VF is set to the value of the most significant bit of VX before the shift
			chip.registers[0xF] = chip.registers[(opcode & 0x0F00) >> 8] >> 7
			chip.registers[(opcode & 0x0F00) >> 8] <<= 1
			chip.pc += 2
	}
}

/**
 * ANNN - Sets I to NNN
 */
func hexA0(chip *Chip8, opcode uint16) {
	chip.index = int(opcode & 0x0FFF)
	chip.pc += 2
}

/**
 * BNNN: Jumps to the address NNN plus V0
 */
func hexB0(chip *Chip8, opcode uint16) {
	chip.pc = int((opcode & 0x0FFF) + uint16(chip.registers[0]))
}

/**
 * CXNN: Sets VX to a random number and NN
 */
func hexC0(chip *Chip8, opcode uint16) {
	chip.registers[(opcode & 0x0F00) >> 8] = byte(int(rand.Int() % 0xFF) & int(opcode & 0x00FF))
	chip.pc += 2
}

/**
 * DXYN: Draws a sprite at coordinate (VX, VY) that has a width of 8 pixels and a height of N pixels. 
 * Each row of 8 pixels is read as bit-coded starting from memory location I; 
 * I value doesn't change after the execution of this instruction. 
 * VF is set to 1 if any screen pixels are flipped from set to unset when the sprite is drawn, 
 * and to 0 if that doesn't happen
 */
func hexD0(chip *Chip8, opcode uint16) {

	x := int(chip.registers[(opcode & 0x0F00) >> 8])
	y := int(chip.registers[(opcode & 0x00F0) >> 4])
	height := int(opcode & 0x000F)
	var pixel uint16

	chip.registers[0xF] = 0

	for yline := 0; yline < height; yline++ {
		pixel = uint16(chip.memory[chip.index + yline])
		for xline := 0; xline < 8; xline++ {
			if pixel & (0x80 >> uint(xline)) != 0 {
				if chip.gfx[(x + xline + ((y + yline) * 64))] == 1 {
					chip.registers[0xF] = 1
				}
				chip.gfx[x + xline + ((y + yline) * 64)] ^= 1
			}
		}
	}

	chip.DrawFlag = true
	chip.pc += 2
}

func hexE0(chip *Chip8, opcode uint16) {

	switch opcode & 0x00FF {
		case 0x009E: // EX9E: Skips the next instruction if the key stored in VX is pressed
			if chip.key[chip.registers[(opcode & 0x0F00) >> 8]] != 0 {
				chip.pc += 4
			} else {
				chip.pc += 2
			}

		case 0x00A1: // EXA1: Skips the next instruction if the key stored in VX isn't pressed
			if chip.key[chip.registers[(opcode & 0x0F00) >> 8]] == 0 {
				chip.pc += 4
			} else {
				chip.pc += 2
			}
	}

}

func hexF0(chip *Chip8, opcode uint16) {
	switch opcode & 0x00FF {
		case 0x0007: // FX07: Sets VX to the value of the delay timer
			chip.registers[(opcode & 0x0F00) >> 8] = chip.delayTimer;
			chip.pc += 2

		case 0x000A: // FX0A: A key press is awaited, and then stored in VX	
			keyPress := false

			for i := 0; i < 16; i++ {
				if chip.key[i] != 0 {
					chip.registers[(opcode & 0x0F00) >> 8] = byte(i)
					keyPress = true
				}
			}

			// If no key is pressed, skip this cycle and try again
			if !keyPress {
				return
			}

			chip.pc += 2

		case 0x0015:
			chip.delayTimer = chip.registers[(opcode & 0x0F00) >> 8]
			chip.pc += 2

		case 0x0018:
			chip.soundTimer = chip.registers[(opcode & 0x0F00) >> 8]
			chip.pc += 2

		case 0x001E:
			if chip.index + int(chip.registers[(opcode & 0x0F00) >> 8]) > 0x0FFF {
				chip.registers[0xF] = 1
			} else {
				chip.registers[0xF] = 0
			}

			chip.index += int(chip.registers[(opcode & 0x0F00) >> 8])
			chip.pc += 2

		case 0x0029:
			chip.index = int(chip.registers[(opcode & 0x0F00) >> 8] * 0x5)
			chip.pc += 2

		case 0x0033:
			chip.memory[chip.index] = chip.registers[(opcode & 0x0F00) >> 8] / 100
			chip.memory[chip.index + 1] = (chip.registers[(opcode & 0x0F00) >> 8] / 10) % 10
			chip.memory[chip.index + 2] = (chip.registers[(opcode & 0x0F00) >> 8] % 100) % 10
			chip.pc += 2

		case 0x0055: // FX55: Stores V0 to VX in memory starting at address I
			for i := 0; i <= int(opcode & 0x0F00) >> 8; i++ {
				chip.memory[chip.index + i] = chip.registers[i]
			}

			// On the original interpreter, when the operation is done, I = I + X + 1.
			chip.index += int(((opcode & 0x0F00) >> 8) + 1)
			chip.pc += 2

		case 0x0065:
			for i := 0; i <= int(opcode & 0x0F00) >> 8; i++ {
				chip.registers[i] = chip.memory[chip.index + i]
			}

			// On the original interpreter, when the operation is done, I = I + X + 1.
			chip.index += int(((opcode & 0x0F00) >> 8) + 1)
			chip.pc += 2

	}

}