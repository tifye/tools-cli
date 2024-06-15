package linking

import "testing"

func TestCalcCrc8(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		input       []byte
		expectedCrc byte
	}{
		"header crc": {
			input:       []byte{0xFD, 0x0D, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			expectedCrc: 0x63,
		},
		"footer crc (checksum)": {
			input:       []byte{0xFD, 0x0D, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x63, 0x12, 0x5F, 0x4D, 0xB6, 0x16},
			expectedCrc: 0x74,
		},
		"tifapp LinkManager.SetProtocol header": {
			input:       []byte{0xFD, 0x0A, 0x00, 0x25, 0xAE, 0xEF, 0x06, 0x00},
			expectedCrc: 0x74,
		},
		"tifapp LinkManager.SetProtocol checksum": {
			input:       []byte{0xFD, 0x0A, 0x00, 0x25, 0xAE, 0xEF, 0x06, 0x00, 0x74, 0x08, 0x01},
			expectedCrc: 0x28,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			test := test
			t.Parallel()
			crc := calcCrc8(test.input)
			if crc != test.expectedCrc {
				t.Errorf("calcCrc8(%v) returned %v; expected %v", test.input, crc, test.expectedCrc)
			}
		})
	}
}
