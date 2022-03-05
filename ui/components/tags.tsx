import { HStack, Tag } from "@chakra-ui/react";
import { pink } from "../lib/color";

interface TagsProps {
  tags: string[];
}

const Tags = ({ tags, ...props }: TagsProps) => {
  return (
    <HStack>
      {tags.map((t) => (
        <Tag
          key={t}
          colorScheme="blackAlpha"
          variant="solid"
          sx={{ color: "black", backgroundColor: "#0000001a" }}
        >
          {t}
        </Tag>
      ))}
    </HStack>
  );
};

export default Tags;
