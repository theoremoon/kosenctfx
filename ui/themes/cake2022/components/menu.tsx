import Link from "next/link";
import { useRouter } from "next/router";
import React from "react";
import styles from "./menu.module.scss";
import cx from "classnames";
import { MenuProps } from "props/menu";

type AnchorProps = JSX.IntrinsicElements["a"];

const ActiveLink = (props: AnchorProps) => {
  const { asPath } = useRouter();
  let classes = props.className;
  if (asPath.startsWith(props.href || "")) {
    classes = cx(classes, styles.active);
  }
  return (
    <Link href={props.href || ""}>
      <a {...props} className={classes}></a>
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
