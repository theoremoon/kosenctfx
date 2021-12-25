import { ChakraProps, Select } from "@chakra-ui/react";
import { countries } from "country-data";
import { ChangeEventHandler } from "react";
import { bgColor, bgSubColor } from "../lib/color";

type CountrySelectorProps = Omit<
  React.ComponentPropsWithRef<"select">,
  "className"
>;

const CountrySelector = ({ ...props }: CountrySelectorProps) => {
  return (
    <select
      className="bg-transparent border border-pink-600 rounded h-8"
      style={{
        backgroundColor: bgColor,
        backgroundImage: `radial-gradient(${bgSubColor} 1px, transparent 1px)`,
        backgroundSize: `10px 10px`,
      }}
      {...props}
    >
      <option value="" className="bg-transparent"></option>
      {countries.all
        .filter((c) => c.emoji)
        .map((c) => (
          <option key={c.name} value={c.alpha2} className="bg-transparent">
            {c.emoji} {c.name}
          </option>
        ))}
    </select>
  );
};

export default CountrySelector;
