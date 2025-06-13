namespace logistics_frontend.Models
{
    public class SidebarMenuItem
    {
        public string Title { get; set; }
        public string Link { get; set; }
        public List<SidebarMenuItem>? SubItems { get; set; }

        public SidebarMenuItem(string title, string link, List<SidebarMenuItem>? subItems = null)
        {
            Title = title;
            Link = link;
            SubItems = subItems;
        }
    }
}
