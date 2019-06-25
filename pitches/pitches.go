package pitches

const HalfStep = 1.05946309436
const WholeStep = HalfStep * HalfStep
const Octave = 2.0

const (
	First  = 1.05946309436
	Second = First * HalfStep
	Third  = Second * HalfStep
	Fourth = Third * HalfStep
	Fifth  = Fourth * HalfStep
)

const (
	C1  = C2 / Octave
	Cs1 = Cs2 / Octave
	Db1 = Db2 / Octave
	D1  = D2 / Octave
	Ds1 = Ds2 / Octave
	Eb1 = Eb2 / Octave
	E1  = E2 / Octave
	F1  = F2 / Octave
	Fs1 = Fs2 / Octave
	Gb1 = Gb2 / Octave
	G1  = G2 / Octave
	Gs1 = Gs2 / Octave
	Ab1 = Ab2 / Octave
	A1  = A2 / Octave
	As1 = As2 / Octave
	Bb1 = Bb2 / Octave
	B1  = B2 / Octave

	C2  = C3 / Octave
	Cs2 = Cs3 / Octave
	Db2 = Db3 / Octave
	D2  = D3 / Octave
	Ds2 = Ds3 / Octave
	Eb2 = Eb3 / Octave
	E2  = E3 / Octave
	F2  = F3 / Octave
	Fs2 = Fs3 / Octave
	Gb2 = Gb3 / Octave
	G2  = G3 / Octave
	Gs2 = Gs3 / Octave
	Ab2 = Ab3 / Octave
	A2  = A3 / Octave
	As2 = As3 / Octave
	Bb2 = Bb3 / Octave
	B2  = B3 / Octave

	C3  = C4 / Octave
	Cs3 = Cs4 / Octave
	Db3 = Db4 / Octave
	D3  = D4 / Octave
	Ds3 = Ds4 / Octave
	Eb3 = Eb4 / Octave
	E3  = E4 / Octave
	F3  = F4 / Octave
	Fs3 = Fs4 / Octave
	Gb3 = Gb4 / Octave
	G3  = G4 / Octave
	Gs3 = Gs4 / Octave
	Ab3 = Ab4 / Octave
	A3  = A4 / Octave
	As3 = As4 / Octave
	Bb3 = Bb4 / Octave
	B3  = B4 / Octave

	C4  = Db4 / HalfStep
	Cs4 = Db4
	Db4 = D4 / HalfStep
	D4  = Eb4 / HalfStep
	Ds4 = Eb4
	Eb4 = E4 / HalfStep
	E4  = F4 / HalfStep
	F4  = Gb4 / HalfStep
	Fs4 = Gb4
	Gb4 = G4 / HalfStep
	G4  = Ab4 / HalfStep
	Gs4 = Ab4
	Ab4 = A4 / HalfStep
	A4  = 440.0
	As4 = A4 * HalfStep
	Bb4 = As4
	B4  = Bb4 * HalfStep

	C5  = C4 * Octave
	Cs5 = Cs4 * Octave
	Db5 = Db4 * Octave
	D5  = D4 * Octave
	Ds5 = Ds4 * Octave
	Eb5 = Eb4 * Octave
	E5  = E4 * Octave
	F5  = F4 * Octave
	Fs5 = Fs4 * Octave
	Gb5 = Gb4 * Octave
	G5  = G4 * Octave
	Gs5 = Gs4 * Octave
	Ab5 = Ab4 * Octave
	A5  = A4 * Octave
	As5 = As4 * Octave
	Bb5 = Bb4 * Octave
	B5  = B4 * Octave

	C6  = C5 * Octave
	Cs6 = Cs5 * Octave
	Db6 = Db5 * Octave
	D6  = D5 * Octave
	Ds6 = Ds5 * Octave
	Eb6 = Eb5 * Octave
	E6  = E5 * Octave
	F6  = F5 * Octave
	Fs6 = Fs5 * Octave
	Gb6 = Gb5 * Octave
	G6  = G5 * Octave
	Gs6 = Gs5 * Octave
	Ab6 = Ab5 * Octave
	A6  = A5 * Octave
	As6 = As5 * Octave
	Bb6 = Bb5 * Octave
	B6  = B5 * Octave

	C7  = C6 * Octave
	Cs7 = Cs6 * Octave
	Db7 = Db6 * Octave
	D7  = D6 * Octave
	Ds7 = Ds6 * Octave
	Eb7 = Eb6 * Octave
	E7  = E6 * Octave
	F7  = F6 * Octave
	Fs7 = Fs6 * Octave
	Gb7 = Gb6 * Octave
	G7  = G6 * Octave
	Gs7 = Gs6 * Octave
	Ab7 = Ab6 * Octave
	A7  = A6 * Octave
	As7 = As6 * Octave
	Bb7 = Bb6 * Octave
	B7  = B6 * Octave
)
