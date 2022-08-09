import { countries } from "country-data";
import { ChangeEventHandler } from "react";
import styles from "./countrySelector.module.scss";
import cx from "classnames";

type CountrySelectorProps = {
  onChange: ChangeEventHandler<HTMLSelectElement>;
  value: string;
  id: string;
} & React.ComponentProps<"select">;

const CountrySelector = ({
  onChange,
  value,
  id,
  ...props
}: CountrySelectorProps) => {
  const { className, ...prps } = props;

  return (
    <select
      onChange={onChange}
      value={value}
      id={id}
      className={cx(className, styles.countryselector)}
      {...prps}
    >
      <option value=""></option>
      {countries.all
        .filter((c) => c.emoji)
        .map((c) => (
          <option key={c.name} value={c.alpha2}>
            {c.emoji} {c.name}
          </option>
        ))}
    </select>
  );
};

export default CountrySelector;
