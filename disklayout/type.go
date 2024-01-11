package disklayout

type ExtType int

const (
	Ext2 = iota
	Ext3
	Ext4
)

func getExtType(fCompat, fInCompat, fRoCompat uint32) ExtType {
	if probeExt4(fCompat, fInCompat, fRoCompat) {
		return Ext4
	} else if probeExt3(fCompat, fInCompat, fRoCompat) {
		return Ext3
	} else if probeExt2(fCompat, fInCompat, fRoCompat) {
		return Ext2
	}
	return 0
}

func probeExt2(fCompat, fInCompat, fRoCompat uint32) bool {
	if (fCompat & SbHasJournal) != 0 {
		return false
	}

	if (fRoCompat&Ext2FeatureRoCompatUnsupported) != 0 ||
		(fInCompat&Ext2FeatureIncompatUnsupported) != 0 {
		return false
	}

	return true
}

func probeExt3(fCompat, fInCompat, fRoCompat uint32) bool {
	if (fCompat & SbHasJournal) == 0 {
		return false
	}

	if (fRoCompat&Ext3FeatureRoCompatUnsupported) != 0 ||
		(fInCompat&Ext3FeatureIncompatUnsupported) != 0 {
		return false
	}

	return true
}

func probeExt4(fCompat, fInCompat, fRoCompat uint32) bool {
	if (fInCompat & SbJournalDev) != 0 {
		return false
	}

	if (fRoCompat&Ext3FeatureRoCompatUnsupported) != 0 ||
		(fInCompat&Ext3FeatureIncompatUnsupported) != 0 {
		return false
	}

	return true
}
