import { HStack, Text } from "@chakra-ui/react";
import { pink } from "../lib/color";

interface TagsProps {
  tags: string[];
}

const Tags = ({ tags, ...props }: TagsProps) => {
  return (
    <HStack>
      {tags.map((t) => (
        <Text
          key={t}
          sx={{
            borderColor: pink,
            borderWidth: "1px",
            borderStyle: "thick",
            paddingLeft: 1,
            paddingRight: 1,
            borderRadius: "4px",
          }}
        >
          {t}{" "}
        </Text>
      ))}
    </HStack>
  );
};

export default Tags;
