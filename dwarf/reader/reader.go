package reader

import (
	"debug/dwarf"
	"fmt"
)

type Reader struct {
	*dwarf.Reader
	depth int
}

// New returns a reader for the specified dwarf data
func New(data *dwarf.Data) *Reader {
	return &Reader{data.Reader(), 0}
}

// Seek moves the reader to an arbitrary offset
func (reader *Reader) Seek(off dwarf.Offset) {
	reader.depth = 0
	reader.Reader.Seek(off)
}

// SeekToEntry moves the reader to an arbitrary entry
func (reader *Reader) SeekToEntry(entry *dwarf.Entry) error {
	reader.Seek(entry.Offset)
	// Consume the current entry so .Next works as intended
	_, err := reader.Next()
	return err
}

// SeekToFunctionEntry moves the reader to the function that includes the
// specified program counter.
func (reader *Reader) SeekToFunction(pc uint64) (*dwarf.Entry, error) {
	reader.Seek(0)
	for entry, err := reader.Next(); entry != nil; entry, err = reader.Next() {
		if err != nil {
			return nil, err
		}

		if entry.Tag != dwarf.TagSubprogram {
			continue
		}

		lowpc, ok := entry.Val(dwarf.AttrLowpc).(uint64)
		if !ok {
			continue
		}

		highpc, ok := entry.Val(dwarf.AttrHighpc).(uint64)
		if !ok {
			continue
		}

		if lowpc <= pc && highpc >= pc {
			return entry, nil
		}
	}

	return nil, fmt.Errorf("unable to find function context")
}

// SeekToType moves the reader to the type specified by the entry,
// optionally resolving typedefs and pointer types. If the reader is set
// to a struct type the NextMemberVariable call can be used to walk all member data
func (reader *Reader) SeekToType(entry *dwarf.Entry, resolveTypedefs bool, resolvePointerTypes bool) (*dwarf.Entry, error) {
	offset, ok := entry.Val(dwarf.AttrType).(dwarf.Offset)
	if !ok {
		return nil, fmt.Errorf("entry does not have a type attribute")
	}

	// Seek to the first type offset
	reader.Seek(offset)

	// Walk the types to the base
	for typeEntry, err := reader.Next(); typeEntry != nil; typeEntry, err = reader.Next() {
		if err != nil {
			return nil, err
		}

		if typeEntry.Tag == dwarf.TagTypedef && !resolveTypedefs {
			return typeEntry, nil
		}

		if typeEntry.Tag == dwarf.TagPointerType && !resolvePointerTypes {
			return typeEntry, nil
		}

		offset, ok = typeEntry.Val(dwarf.AttrType).(dwarf.Offset)
		if !ok {
			return typeEntry, nil
		}

		reader.Seek(offset)
	}

	return nil, fmt.Errorf("no type entry found")
}

// NextScopeVariable moves the reader to the next debug entry that describes a local variable and returns the entry
func (reader *Reader) NextScopeVariable() (*dwarf.Entry, error) {
	for entry, err := reader.Next(); entry != nil; entry, err = reader.Next() {
		if err != nil {
			return nil, err
		}

		// All scope variables will be at the same depth
		reader.SkipChildren()

		// End of the current depth
		if entry.Tag == 0 {
			break
		}

		if entry.Tag == dwarf.TagVariable || entry.Tag == dwarf.TagFormalParameter {
			return entry, nil
		}
	}

	// No more items
	return nil, nil
}

// NextMememberVariable moves the reader to the next debug entry that describes a member variable and returns the entry
func (reader *Reader) NextMemberVariable() (*dwarf.Entry, error) {
	for entry, err := reader.Next(); entry != nil; entry, err = reader.Next() {
		if err != nil {
			return nil, err
		}

		// All member variables will be at the same depth
		reader.SkipChildren()

		// End of the current depth
		if entry.Tag == 0 {
			break
		}

		if entry.Tag == dwarf.TagMember {
			return entry, nil
		}
	}

	// No more items
	return nil, nil
}

// NextPackage moves the reader to the next debug entry that describes a package variable
// any TagVariable entry that is not inside a sub prgram entry and is marked external is considered a package variable
func (reader *Reader) NextPackageVariable() (*dwarf.Entry, error) {
	for entry, err := reader.Next(); entry != nil; entry, err = reader.Next() {
		if err != nil {
			return nil, err
		}

		if entry.Tag == dwarf.TagVariable {
			ext, ok := entry.Val(dwarf.AttrExternal).(bool)
			if ok && ext {
				return entry, nil
			}
		}

		// Ignore everything inside sub programs
		if entry.Tag == dwarf.TagSubprogram {
			reader.SkipChildren()
		}
	}

	// No more items
	return nil, nil
}
