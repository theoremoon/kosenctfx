export type MenuItem = {
  href: string;
  innerText: string;
};

export interface MenuProps {
  siteName: string;
  leftMenuItems: MenuItem[];
  rightMenuItems: MenuItem[];
}
