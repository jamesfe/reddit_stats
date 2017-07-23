/* Not to be confused with utilities unit tests, these are utilities I use IN the tests. */

package utils

func LoadRawByteArray(filename string, numItems int) [][]byte {
	/* Load an array of so many lines of raw bytes from a file */
	inFileReader, f := GetFileReader(filename)
	defer f()
	var returnValues [][]byte = make([][]byte, numItems)
	for i := 0; i < numItems; i++ {
		if inputBytes, err := inFileReader.ReadBytes('\n'); err != nil {
			log.Errorf("File Error: %s", err) // maybe we are in an IO error?
		} else {
			returnValues[i] = inputBytes
		}
	}
	return returnValues
}
