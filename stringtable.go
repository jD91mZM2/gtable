package gtable;

const (
	CORNER_TOP_LEFT = iota
	CORNER_TOP_RIGHT
	CORNER_BOTTOM_LEFT
	CORNER_BOTTOM_RIGHT
	CORNER_MIDDLE_HORIZONTAL
	CORNER_MIDDLE_VERTICAL
)
var CORNER_PLUS = func(i int) rune{
	return '+';
};
var CORNER_ROUND = func(i int) rune{
	switch(i){
		default:
			return '.';
		case CORNER_MIDDLE_HORIZONTAL:
			fallthrough;
		case CORNER_TOP_RIGHT:
			return '+';
		case CORNER_BOTTOM_LEFT:
			fallthrough;
		case CORNER_BOTTOM_RIGHT:
			return '\'';
	}
};

type stringTable struct{
	Corner func(i int) rune
	Header bool
	rows [][]*TableItem
}
func NewStringTable() stringTable{
	return stringTable{
		Corner: CORNER_PLUS,
		Header: true,
		rows: [][]*TableItem{ make([]*TableItem, 0) },
	};
}

func (st *stringTable) AddItems(items ...*TableItem){
	index := len(st.rows) - 1;
	col := st.rows[index];

	for _, item := range items{
		col = append(col, item);
	}
	st.rows[index] = col;
}
func (st *stringTable) AddStrings(items ...string){
	tItems := make([]*TableItem, len(items));
	for i, item := range items{
		tItem := NewItem(item);
		tItems[i] = &tItem;
	}
	st.AddItems(tItems...);
}
func (st *stringTable) AddRow(){
	st.rows = append(st.rows, make([]*TableItem, 0))
}
func (st *stringTable) Get(row, col int) *TableItem{
	return st.rows[row][col];
}

func (st *stringTable) Columns() int{
	columns := 0;

	for _, row := range st.rows{
		columns = max(len(row), columns);
	}
	return columns;
}

func max(i1, i2 int) int{
	if(i1 >= i2){
		return i1;
	} else {
		return i2;
	}
}

func (st *stringTable) Rows() [][]*TableItem{
	var arr = make([][]*TableItem, len(st.rows));
	for i := range st.rows{
		arr[i] = make([]*TableItem, len(st.rows[i]));
		copy(arr[i], st.rows[i])
	}

	return arr;
}
func (st *stringTable) Each(handler func(*TableItem)){
	for _, row := range st.rows{
		for _, col := range row{
			handler(col);
		}
	}
}

func (st *stringTable) String() string{
	s := "";
	n := "\n";

	rows := st.Rows();

	last := len(rows) - 1;
	for len(rows) > 0 && len(rows[last]) <= 0{
		rows = append(rows[:last], rows[last:]...);
	}

	if(len(rows) <= 0){
		return "";
	}

	columns := 0;
	for _, row := range rows{
		columns = max(columns, len(row));
	}

	if(columns <= 0){
		return "";
	}

	lengths := make([]int, 0);
	for c := 0; c < columns; c++{
		length := 0;
		for r := 0; r < len(rows); r++{
			for len(rows[r]) < columns{
				item := NewItem("");
				rows[r] = append(rows[r], &item);
			}
			col := rows[r][c];

			length = max(length, col.Size());
		}
		lengths = append(lengths, length);
	}

	frame := "";

	first := true;
	for _, length := range lengths{
		if(first){
			first = false;
		} else {
			frame += string(st.Corner(CORNER_MIDDLE_HORIZONTAL));
		}

		for i := 0; i < length; i++ {
			frame += "-";
		}
	}

	s += string(st.Corner(CORNER_TOP_LEFT)) + frame + string(st.Corner(CORNER_TOP_RIGHT)) + n;

	first = true;
	for _, row := range rows{
		for i, col := range row{
			col.Width = lengths[i] - (col.PaddingLeft + col.PaddingRight);

			s += "|";
			s += col.String();
		}
		s += "|" + n;

		if(first){
			first = false;

			if(st.Header){
				s += string(st.Corner(CORNER_MIDDLE_VERTICAL)) + frame + string(st.Corner(CORNER_MIDDLE_VERTICAL)) + n;
			}
		}
	}

	s += string(st.Corner(CORNER_BOTTOM_LEFT)) + frame + string(st.Corner(CORNER_BOTTOM_RIGHT));

	return s;
}