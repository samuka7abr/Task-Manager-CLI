package main 

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

type Priority string

const (
	Alta Priority = "alta"
	Media Priority = "media"
	Baixa Priority = "baixa"	
)

type Task struct {
	ID	  int	    `json:"id"`
	Name	  string    `json:"name`
	Priority  Priority  `json:"priority"`
	DueDate   time.Time `json:"due_date`
	CreatedAt time.Time `json:"created_at`
}

type Store struct {
	Tasks  []Task `json:"tasks"`
	NextID int    `json:"next_id"`
}

var (
	storePath string
	db	  Store
)

var rootCmd = &cobra.Command{
	Use:               "taskman",
	Short:             "Gerenciador de tarefas em CLI",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error { return loadStore() },
	PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
		return saveStore()
	},
}

var addCmd = &cobra.Command{
	Use:   "add <nome> [prioridade] [dia] [mes] [ano]",
	Short: "Adicionar tarefa",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		nome := args[0]
		p := Media
		if len(args) >= 2 && args[1] != "" {
			a := strings.ToLower(args[1])
			switch a {
			case "alta":
				p = Alta
			case "media", "média":
				p = Media
			case "baixa":
				p = Baixa
			default:
				return fmt.Errorf("prioridade inválida: use alta, media, baixa")
			}
		}
		now := time.Now()
		d := now.Day()
		m := int(now.Month())
		y := now.Year()
		if len(args) >= 3 && args[2] != "" {
			v, err := strconv.Atoi(args[2])
			if err != nil || v < 1 || v > 31 {
				return fmt.Errorf("dia inválido")
			}
			d = v
		}
		if len(args) >= 4 && args[3] != "" {
			v, err := strconv.Atoi(args[3])
			if err != nil || v < 1 || v > 12 {
				return fmt.Errorf("mês inválido")
			}
			m = v
		}
		if len(args) >= 5 && args[4] != "" {
			v, err := strconv.Atoi(args[4])
			if err != nil || v < 1 {
				return fmt.Errorf("ano inválido")
			}
			y = v
		}
		date, err := safeDate(y, m, d)
		if err != nil {
			return err
		}
		id := db.NextID
		if id == 0 {
			id = 1
		}
		db.Tasks = append(db.Tasks, Task{
			ID:        id,
			Name:      nome,
			Priority:  p,
			DueDate:   date,
			CreatedAt: time.Now(),
		})
		db.NextID = id + 1
		fmt.Printf("Adicionada #%d: %s [%s] %s\n", id, nome, strings.ToUpper(string(p)), date.Format("02/01/2006"))
		return nil
	},
}

var delCmd = &cobra.Command{
	Use:   "del <id>",
	Short: "Deletar tarefa por ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("id inválido")
		}
		i := indexByID(id)
		if i == -1 {
			return fmt.Errorf("tarefa #%d não encontrada", id)
		}
		t := db.Tasks[i]
		db.Tasks = append(db.Tasks[:i], db.Tasks[i+1:]...)
		fmt.Printf("Deletada #%d: %s %s\n", t.ID, t.Name, t.DueDate.Format("02/01/2006"))
		return nil
	},
}

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Editar tarefa (fluxo interativo)",
	RunE: func(cmd *cobra.Command, args []string) error {
		r := bufio.NewReader(os.Stdin)
		fmt.Print("Data da tarefa (dd [mm] [aaaa], vazio = hoje): ")
		ds, _ := r.ReadString('\n')
		ds = strings.TrimSpace(ds)
		dt, err := parseFlexibleDate(ds)
		if err != nil {
			return err
		}
		fmt.Print("Nome exato da tarefa nessa data: ")
		ns, _ := r.ReadString('\n')
		ns = strings.TrimSpace(ns)
		matches := findByDateAndName(dt, ns)
		if len(matches) == 0 {
			return fmt.Errorf("nenhuma tarefa encontrada em %s com nome '%s'", dt.Format("02/01/2006"), ns)
		}
		var t *Task
		if len(matches) == 1 {
			t = matches[0]
		} else {
			fmt.Println("Múltiplas tarefas encontradas:")
			for _, x := range matches {
				fmt.Printf("ID:%d %s [%s]\n", x.ID, x.Name, x.DueDate.Format("02/01/2006"))
			}
			fmt.Print("Digite o ID para editar: ")
			ids, _ := r.ReadString('\n')
			id, err := strconv.Atoi(strings.TrimSpace(ids))
			if err != nil {
				return fmt.Errorf("id inválido")
			}
			idx := indexByID(id)
			if idx == -1 {
				return fmt.Errorf("id não encontrado")
			}
			t = &db.Tasks[idx]
		}
		fmt.Printf("Novo nome (enter mantém: \"%s\"): ", t.Name)
		nn, _ := r.ReadString('\n')
		nn = strings.TrimSpace(nn)
		if nn != "" {
			t.Name = nn
		}
		fmt.Printf("Nova prioridade [alta|media|baixa] (enter mantém: %s): ", t.Priority)
		pp, _ := r.ReadString('\n')
		pp = strings.TrimSpace(strings.ToLower(pp))
		switch pp {
		case "":
		case "alta":
			t.Priority = Alta
		case "media", "média":
			t.Priority = Media
		case "baixa":
			t.Priority = Baixa
		default:
			return fmt.Errorf("prioridade inválida")
		}
		fmt.Printf("Nova data (dd [mm] [aaaa], enter mantém: %s): ", t.DueDate.Format("02/01/2006"))
		nd, _ := r.ReadString('\n')
		nd = strings.TrimSpace(nd)
		if nd != "" {
			dd, err := parseFlexibleDate(nd)
			if err != nil {
				return err
			}
			t.DueDate = dd
		}
		fmt.Printf("Editada #%d: %s [%s] %s\n", t.ID, t.Name, strings.ToUpper(string(t.Priority)), t.DueDate.Format("02/01/2006"))
		return nil
	},
}

var dayCmd = &cobra.Command{
	Use:   "day [dia] [mes] [ano]",
	Short: "Listar tarefas do dia",
	RunE: func(cmd *cobra.Command, args []string) error {
		dt, err := dateFromArgs(args)
		if err != nil {
			return err
		}
		items := tasksOnDate(dt)
		printTasksTable(items, fmt.Sprintf("Tarefas de %s", dt.Format("02/01/2006")))
		return nil
	},
}

var weekCmd = &cobra.Command{
	Use:   "week [dia] [mes] [ano]",
	Short: "Listar tarefas da semana",
	RunE: func(cmd *cobra.Command, args []string) error {
		ref, err := dateFromArgs(args)
		if err != nil {
			return err
		}
		start := monday(ref)
		end := start.AddDate(0, 0, 7)
		var items []Task
		for _, t := range db.Tasks {
			if !t.DueDate.Before(start) && t.DueDate.Before(end) {
				items = append(items, t)
			}
		}
		printTasksTable(items, fmt.Sprintf("Semana %s a %s", start.Format("02/01"), end.AddDate(0, 0, -1).Format("02/01")))
		return nil
	},
}

var calCmd = &cobra.Command{
	Use:   "cal [mes] [ano]",
	Short: "Calendário com mapa de calor",
	RunE: func(cmd *cobra.Command, args []string) error {
		now := time.Now()
		m := int(now.Month())
		y := now.Year()
		if len(args) >= 1 && args[0] != "" {
			v, err := strconv.Atoi(args[0])
			if err != nil || v < 1 || v > 12 {
				return fmt.Errorf("mês inválido")
			}
			m = v
		}
		if len(args) >= 2 && args[1] != "" {
			v, err := strconv.Atoi(args[1])
			if err != nil || v < 1 {
				return fmt.Errorf("ano inválido")
			}
			y = v
		}
		printCalendar(y, m)
		return nil
	},
}

func main() {
	rootCmd.AddCommand(addCmd, delCmd, editCmd, dayCmd, weekCmd, calCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func loadStore() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	storePath = filepath.Join(home, ".taskman_cli.json")
	_, err = os.Stat(storePath)
	if os.IsNotExist(err) {
		db = Store{NextID: 1}
		return nil
	}
	b, err := os.ReadFile(storePath)
	if err != nil {
		return err
	}
	if len(b) == 0 {
		db = Store{NextID: 1}
		return nil
	}
	if err := json.Unmarshal(b, &db); err != nil {
		return err
	}
	if db.NextID == 0 {
		maxID := 0
		for _, t := range db.Tasks {
			if t.ID > maxID {
				maxID = t.ID
			}
		}
		db.NextID = maxID + 1
	}
	return nil
}

func saveStore() error {
	b, err := json.MarshalIndent(db, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(storePath, b, 0644)
}

func safeDate(y, m, d int) (time.Time, error) {
	loc := time.Now().Location()
	t := time.Date(y, time.Month(m), d, 0, 0, 0, 0, loc)
	if t.Month() != time.Month(m) || t.Day() != d || t.Year() != y {
		return time.Time{}, fmt.Errorf("data inválida")
	}
	return t, nil
}

func parseFlexibleDate(s string) (time.Time, error) {
	now := time.Now()
	if strings.TrimSpace(s) == "" {
		return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()), nil
	}
	parts := strings.Fields(s)
	d, m, y := now.Day(), int(now.Month()), now.Year()
	if len(parts) >= 1 {
		v, err := strconv.Atoi(parts[0])
		if err != nil {
			return time.Time{}, fmt.Errorf("dia inválido")
		}
		d = v
	}
	if len(parts) >= 2 {
		v, err := strconv.Atoi(parts[1])
		if err != nil || v < 1 || v > 12 {
			return time.Time{}, fmt.Errorf("mês inválido")
		}
		m = v
	}
	if len(parts) >= 3 {
		v, err := strconv.Atoi(parts[2])
		if err != nil || v < 1 {
			return time.Time{}, fmt.Errorf("ano inválido")
		}
		y = v
	}
	return safeDate(y, m, d)
}

func dateFromArgs(args []string) (time.Time, error) {
	if len(args) == 0 {
		return parseFlexibleDate("")
	}
	return parseFlexibleDate(strings.Join(args, " "))
}

func tasksOnDate(dt time.Time) []Task {
	var out []Task
	for _, t := range db.Tasks {
		if sameDay(t.DueDate, dt) {
			out = append(out, t)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].DueDate.Before(out[j].DueDate) || out[i].ID < out[j].ID })
	return out
}

func sameDay(a, b time.Time) bool {
	ay, am, ad := a.Date()
	by, bm, bd := b.Date()
	return ay == by && am == bm && ad == bd
}

func indexByID(id int) int {
	for i, t := range db.Tasks {
		if t.ID == id {
			return i
		}
	}
	return -1
}

func findByDateAndName(dt time.Time, name string) []*Task {
	var out []*Task
	for i := range db.Tasks {
		if sameDay(db.Tasks[i].DueDate, dt) && db.Tasks[i].Name == name {
			out = append(out, &db.Tasks[i])
		}
	}
	return out
}

func colorPriority(p Priority, s string) string {
	switch p {
	case Alta:
		return "\033[31m" + s + "\033[0m"
	case Media:
		return "\033[33m" + s + "\033[0m"
	default:
		return "\033[32m" + s + "\033[0m"
	}
}

func printTasksTable(items []Task, title string) {
	if len(items) == 0 {
		fmt.Println(title + ": nenhum item.")
		return
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].DueDate.Equal(items[j].DueDate) {
			return items[i].ID < items[j].ID
		}
		return items[i].DueDate.Before(items[j].DueDate)
	})
	fmt.Println(title)
	fmt.Println(strings.Repeat("-", 72))
	fmt.Printf("%-4s %-10s %-12s %-42s\n", "ID", "PRIORIDADE", "DATA", "NOME")
	for _, t := range items {
		pr := strings.ToUpper(string(t.Priority))
		fmt.Printf("%-4d %-10s %-12s %-42s\n", t.ID, colorPriority(t.Priority, pr), t.DueDate.Format("02/01/2006"), truncate(t.Name, 42))
	}
}

func truncate(s string, max int) string {
	if len([]rune(s)) <= max {
		return s
	}
	r := []rune(s)
	return string(r[:max-1]) + "…"
}

func monday(t time.Time) time.Time {
	wd := int(t.Weekday())
	if wd == 0 {
		wd = 7
	}
	return time.Date(t.Year(), t.Month(), t.Day()-(wd-1), 0, 0, 0, 0, t.Location())
}

func monthNamePT(m int) string {
	n := []string{"", "janeiro", "fevereiro", "março", "abril", "maio", "junho", "julho", "agosto", "setembro", "outubro", "novembro", "dezembro"}
	return n[m]
}

func printCalendar(year, month int) {
	loc := time.Now().Location()
	first := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, loc)
	next := first.AddDate(0, 1, 0)
	counts := map[int]int{}
	maxc := 0
	for _, t := range db.Tasks {
		if t.DueDate.Before(first) || !t.DueDate.Before(next) {
			continue
		}
		_, _, d := t.DueDate.Date()
		counts[d]++
		if counts[d] > maxc {
			maxc = counts[d]
		}
	}
	fmt.Printf("%s de %d\n", strings.Title(monthNamePT(month)), year)
	fmt.Println("Seg Ter Qua Qui Sex Sab Dom")
	startOffset := weekdayOffsetMonday(first)
	day := 1
	today := time.Now()
	for i := 0; i < startOffset; i++ {
		fmt.Printf("    ")
	}
	for day <= daysInMonth(year, month) {
		level := 0
		if maxc > 0 {
			level = int(float64(counts[day]) / float64(maxc) * 4.0)
		}
		cell := heatBG(level) + fmt.Sprintf("%2d", day) + "\033[0m"
		if sameDay(first.AddDate(0, 0, day-1), today) {
			cell = "\033[1m" + cell + "\033[0m"
		}
		fmt.Printf(" %s", cell)
		if (day+startOffset)%7 == 0 {
			fmt.Println()
		}
		day++
	}
	if (day+startOffset-1)%7 != 0 {
		fmt.Println()
	}
	fmt.Println("  legendas: intensidade baseada em quantidade de tarefas")
}

func daysInMonth(y, m int) int {
	t := time.Date(y, time.Month(m), 1, 0, 0, 0, 0, time.Now().Location())
	return t.AddDate(0, 1, -1).Day()
}

func weekdayOffsetMonday(t time.Time) int {
	wd := int(t.Weekday())
	if wd == 0 {
		wd = 7
	}
	return wd - 1
}

func heatBG(level int) string {
	switch level {
	case 0:
		return "\033[48;5;236m\033[38;5;250m"
	case 1:
		return "\033[48;5;242m\033[38;5;255m"
	case 2:
		return "\033[48;5;148m\033[38;5;16m"
	case 3:
		return "\033[48;5;112m\033[38;5;16m"
	default:
		return "\033[48;5;46m\033[38;5;16m"
	}
}
