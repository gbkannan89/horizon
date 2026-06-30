package engine

// Category represents the notification category.
type Category string
const (
	CatCritical       Category = "Critical"
	CatWarning        Category = "Warning"
	CatReminder       Category = "Reminder"
	CatRecommendation Category = "Recommendation"
	CatAchievement    Category = "Achievement"
	CatInsight        Category = "Insight"
	CatGoal           Category = "Goal"
	CatRisk           Category = "Risk"
	CatPlanning       Category = "Planning"
	CatPortfolio      Category = "Portfolio"
	CatHealth         Category = "Health"
	CatSystem         Category = "System"
)

// Priority represents the delivery priority level.
type Priority string
const (P1Critical Priority = "P1"; P2Urgent Priority = "P2"; P3Important Priority = "P3"; P4Info Priority = "P4")

// NotificationState represents the lifecycle state.
type NotificationState string
const (NSGenerated NotificationState = "Generated"; NSQueued NotificationState = "Queued"; NSDelivered NotificationState = "Delivered"; NSViewed NotificationState = "Viewed"; NSActedUpon NotificationState = "ActedUpon"; NSDismissed NotificationState = "Dismissed"; NSExpired NotificationState = "Expired"; NSArchivedSt NotificationState = "Archived")

// NotificationCard represents a single notification.
type NotificationCard struct {
	NotifID       string            `json:"notif_id"`
	Category      Category          `json:"category"`
	Priority      Priority          `json:"priority"`
	Title         string            `json:"title"`
	Summary       string            `json:"summary"`
	ActionLabel   string            `json:"action_label"`
	ActionRoute   string            `json:"action_route"`
	State         NotificationState `json:"state"`
	GeneratedAt   string            `json:"generated_at"`
	ExpiresAt     string            `json:"expires_at"`
}

// NotificationCenter is the complete notification view.
type NotificationCenter struct {
	UnreadCount int                `json:"unread_count"`
	P1Count     int                `json:"p1_count"`
	P2Count     int                `json:"p2_count"`
	P3Count     int                `json:"p3_count"`
	P4Count     int                `json:"p4_count"`
	Items       []NotificationCard `json:"items"`
}

// Preference represents a per-category notification preference.
type Preference struct {
	Category       Category `json:"category"`
	InApp          bool     `json:"in_app"`
	Push           bool     `json:"push"`
	Email          bool     `json:"email"`
	Digest         string   `json:"digest"` // none, daily, weekly
	MinPriority    string   `json:"min_priority"`
	Enabled        bool     `json:"enabled"`
}

// Inputs for the notifications experience.
type Inputs struct {
	UserID      string            `json:"user_id"`
	Notifications []NotifInput    `json:"notifications"`
	Preferences []PreferenceInput `json:"preferences"`
}

// NotifInput represents raw notification data.
type NotifInput struct {
	NotifID     string `json:"notif_id"`
	Category    string `json:"category"`
	Priority    string `json:"priority"`
	Title       string `json:"title"`
	Summary     string `json:"summary"`
	ActionLabel string `json:"action_label"`
	ActionRoute string `json:"action_route"`
	State       string `json:"state"`
	GeneratedAt string `json:"generated_at"`
	ExpiresAt   string `json:"expires_at"`
}

// PreferenceInput represents raw preference data.
type PreferenceInput struct {
	Category    string `json:"category"`
	InApp       bool   `json:"in_app"`
	Push        bool   `json:"push"`
	Email       bool   `json:"email"`
	Digest      string `json:"digest"`
	MinPriority string `json:"min_priority"`
	Enabled     bool   `json:"enabled"`
}
