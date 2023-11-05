import Link, { LinkProps } from "next/link";
import { useRouter } from "next/router";
import React, { ReactNode } from "react";
import styles from "./menu.module.scss";
import cx from "classnames";
import { MenuProps } from "props/menu";

type ActiveLinkProps = LinkProps & { children: ReactNode };

const ActiveLink = (props: ActiveLinkProps) => {
  const { asPath } = useRouter();
  let classes = "";
  if (asPath.startsWith(props.href.toString() || "")) {
    classes = cx(classes, styles.active);
  }
  return (
    <Link href={props.href || ""} className={classes}>
      {props.children}
    </Link>
  );
};

const Menu = ({ leftMenuItems, rightMenuItems }: MenuProps) => {
  return (
    <nav className={styles.nav}>
      <ul>
        {[...leftMenuItems, ...rightMenuItems].map((item) => (
          <li key={item.innerText}>
            <ActiveLink href={item.href}>{item.innerText}</ActiveLink>
          </li>
        ))}
      </ul>
    </nav>
  );
};
export default Menu;
