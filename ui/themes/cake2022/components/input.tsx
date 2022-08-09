import React, { ComponentPropsWithoutRef } from "react";
import styles from "./input.module.scss";
import cx from "classnames";

type InputProps = ComponentPropsWithoutRef<"input">;

// eslint-disable-next-line react/display-name
const Input = React.forwardRef<HTMLInputElement, InputProps>((props, ref) => {
  const { className, ...ps } = props;
  const cs = cx(className || "", styles.text);
  return <input ref={ref} className={cs} {...ps} />;
});

export default Input;
