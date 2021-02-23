package timeshift

import (
	"fmt"
	"regexp"
	"strconv"
	"sync"
	"time"
)

//----------------------------------------------------------------------------------------------------------------------------//

type (
	// TimeShift --
	TimeShift struct {
		empty   bool
		year    partDef
		month   partDef
		day     partDef
		week    partDef
		weekday partDef
		hour    partDef
		minute  partDef
		second  partDef
	}

	partDef struct {
		active    bool
		val       int
		absolute  bool
		fromBegin bool // for week only
		fromEnd   bool // for day and week only
	}
)

var (
	partExpression  = `(?:\s*)([YMDWwhms])([\^\$]?)([+-]?)(\d+)(?:\s*)`
	checkExpression = fmt.Sprintf(`^(%s)+$`, partExpression)

	checkRE    = regexp.MustCompile(checkExpression)
	splitRE    = regexp.MustCompile(partExpression)
	cacheMutex = new(sync.RWMutex)
	cache      = map[string]*TimeShift{}
)

const (
	partSrc     = 0
	partName    = 1
	partOptions = 2
	partSign    = 3
	partVal     = 4
)

//----------------------------------------------------------------------------------------------------------------------------//

// New --
func New(pattern string, cached bool) (ts *TimeShift, err error) {
	if cached {
		cacheMutex.RLock()
		ts = cache[pattern]
		cacheMutex.RUnlock()

		if ts != nil {
			return
		}
	}

	ts = &TimeShift{empty: false}

	defer func() {
		if err != nil {
			ts = nil
		} else if cached {
			cacheMutex.Lock()
			cache[pattern] = ts
			cacheMutex.Unlock()
		}
	}()

	if pattern == "" {
		ts.empty = true
		return
	}

	// Common check
	parts := checkRE.FindAllStringSubmatch(pattern, -1)
	if len(parts) == 0 {
		err = fmt.Errorf(`Illegal pattern "%s"`, pattern)
		return
	}

	// Split to parts
	parts = splitRE.FindAllStringSubmatch(pattern, -1)

	// Parts sequence pattern
	partNames := []byte("YMDWwhms!")
	nameIdx := 0

	for _, part := range parts {
		name := part[partName]

		// Checking the sequence of parts and their uniqueness
		for ; nameIdx < len(partNames); nameIdx++ {
			if byte(name[0]) == partNames[nameIdx] {
				nameIdx++
				break
			}
		}
		if nameIdx >= len(partNames) {
			err = fmt.Errorf(`Wrong sequence of parts in "%s" (about %s), expected "%s"`, pattern, part[partSrc], partNames[:len(partNames)-1])
			return
		}

		v, _ := strconv.ParseInt(part[partVal], 10, 32)

		pDf := partDef{
			active:    true,
			val:       int(v),
			absolute:  true,
			fromBegin: false,
			fromEnd:   false,
		}

		switch part[partSign] {
		case "+":
			pDf.absolute = false
		case "-":
			pDf.absolute = false
			pDf.val = -pDf.val
		}

		for _, c := range part[partOptions] {
			switch c {
			case '^':
				switch name {
				case "W":
					pDf.fromBegin = true
				default:
					err = fmt.Errorf(`Illegal option "%c" in the "%s"`, c, part[partSrc])
					return
				}
			case '$':
				switch name {
				case "D", "W":
					pDf.fromEnd = true
				default:
					err = fmt.Errorf(`Illegal option "%c" in the "%s"`, c, part[partSrc])
					return
				}
			}
		}

		if (pDf.fromBegin || pDf.fromEnd) && !pDf.absolute {
			err = fmt.Errorf(`"^" and "$" can not be used with relative ("+" or "-") values in the "%s"`, part[partSrc])
			return
		}

		if pDf.fromBegin && pDf.fromEnd {
			err = fmt.Errorf(`"^" and "$" can not be used simultaneously in the "%s"`, part[partSrc])
			return
		}

		switch name {
		case "Y":
			ts.year = pDf

		case "M":
			if pDf.absolute && pDf.val == 0 {
				err = fmt.Errorf(`Illegal month in the "%s"`, part[partSrc])
				return
			}
			ts.month = pDf

		case "D":
			if pDf.active {
				if pDf.absolute && pDf.val == 0 {
					err = fmt.Errorf(`Illegal day in the "%s"`, part[partSrc])
					return
				}
			}
			ts.day = pDf

		case "W":
			if pDf.active {
				if pDf.fromBegin || pDf.fromEnd {
					if pDf.val == 0 {
						err = fmt.Errorf(`Illegal relative week in the "%s"`, part[partSrc])
						return
					}
				} else if pDf.absolute {
					if pDf.val == 0 {
						err = fmt.Errorf(`Illegal absolute week in the "%s"`, part[partSrc])
						return
					}
				}
			}
			ts.week = pDf

		case "w":
			// 0 - Sunday
			if pDf.val < 0 || pDf.val > 6 {
				err = fmt.Errorf(`Illegal weekday in the "%s"`, part[partSrc])
				return
			}
			ts.weekday = pDf

		case "h":
			ts.hour = pDf

		case "m":
			ts.minute = pDf

		case "s":
			ts.second = pDf
		}
	}

	return
}

//----------------------------------------------------------------------------------------------------------------------------//

// Exec --
func (ts *TimeShift) Exec(t time.Time) (result time.Time) {
	if ts.empty {
		result = t
		return
	}

	proc := func(df *partDef, v *int) {
		if !df.active {
			return
		}

		if df.fromEnd {
			return // save the old value temporarily
		}

		if df.absolute {
			*v = df.val
			return
		}

		*v += df.val
		return
	}

	hour, minute, second := t.Clock()
	year, m, day := t.Date()
	month := int(m)

	proc(&ts.hour, &hour)
	proc(&ts.minute, &minute)
	proc(&ts.second, &second)

	proc(&ts.year, &year)
	proc(&ts.month, &month)
	proc(&ts.day, &day)

	result = time.Date(year, time.Month(month), day, hour, minute, second, 0, t.Location())

	if ts.day.fromEnd {
		result = result.AddDate(0, 1, -day-ts.day.val+1)
	}

	if ts.week.active {
		df := ts.week

		var wd int
		if ts.weekday.active {
			wd = ts.weekday.val
		} else {
			wd = int(result.Weekday())
		}

		if df.fromBegin {
			// from the begin of the month
			result = result.AddDate(0, 0, -result.Day()+1) // begin of the month

			shift := wd - int(result.Weekday())
			if shift < 0 {
				shift += 7
			}
			shift += (df.val - 1) * 7

			result = result.AddDate(0, 0, shift)
			return // weekday already taken
		}

		if df.fromEnd {
			// from the end of the month
			result = result.AddDate(0, 1, -result.Day()) // end of the month

			shift := wd - int(result.Weekday())
			if shift > 0 {
				shift -= 7
			}
			shift -= (df.val - 1) * 7

			result = result.AddDate(0, 0, shift)
			return // weekday already taken
		}

		if df.absolute {
			// from begin of the year
			result = result.AddDate(0, 0, -result.YearDay()+1) // 1 Jan

			shift := wd - int(result.Weekday())
			if shift < 0 {
				shift += 7
			}
			shift += (df.val - 1) * 7

			result = result.AddDate(0, 0, shift)
			return // weekday already taken
		}

		// relative the result date - simple shift and continue to weekday
		result = result.AddDate(0, 0, df.val*7)
	}

	if ts.weekday.active {
		shift := ts.weekday.val - int(result.Weekday())
		result = result.AddDate(0, 0, shift)
		return
	}

	return
}

//----------------------------------------------------------------------------------------------------------------------------//
