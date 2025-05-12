export interface Announcement {
  text: string;
  date: string;
  isNew: boolean;
}

const announcements: Announcement[] = [
  {
    text: "Welcome to Pagoda on Neocities. The site is currently under construction and as we progress and add more features and content, this announcement section will get updated alongside it. Stay tuned for more updates!",
    date: "2025-05-13 02:00",
    isNew: false,
  },
  {
    text: "Announcement section and Districts are now added!",
    date: "2025-05-13 02:15",
    isNew: true,
  },
];

export const getAnnouncements = (): Announcement[] => {
  return announcements.sort((a, b) => {
    return new Date(b.date).getTime() - new Date(a.date).getTime();
  });
};
