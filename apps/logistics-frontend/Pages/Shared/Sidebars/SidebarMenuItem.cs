namespace logistics_frontend.Models
{
    public class SidebarMenuItem
    {
        public string Title { get; set; }
        public string Link { get; set; }
        public string Icon { get; set; }
        public List<SidebarMenuItem>? SubItems { get; set; }
        public int? BadgeCount { get; set; }

        public SidebarMenuItem(string title, string link, string icon, List<SidebarMenuItem>? subItems = null, int? badgeCount = null)
        {
            Title = title;
            Link = link;
            Icon = icon;
            SubItems = subItems;
            BadgeCount = badgeCount;
        }
    }
}
