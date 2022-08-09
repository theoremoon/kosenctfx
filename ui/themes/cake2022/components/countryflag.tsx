import { lookup } from "country-data";

type CountryFlagProps = {
  country: string;
};

const CountryFlag = ({ country }: CountryFlagProps) => {
  const c = lookup.countries({ alpha2: country })[0];
  if (!c) {
    return <></>;
  }
  return <span title={c.name}>{c.emoji}</span>;
};

export default CountryFlag;
