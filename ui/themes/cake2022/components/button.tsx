import cx from "classnames";
import styles from "./button.module.scss";

type ButtonProps = React.ComponentProps<"button">;

const Button = ({ ...props }: ButtonProps) => {
  props.className = cx(props.className, styles.button);

  return <button {...props} />;
};

export default Button;
