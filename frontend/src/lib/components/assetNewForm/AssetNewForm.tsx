import { newFolderPost, patchAsset } from "@/lib/fetchers/getAsset";
import { Asset, AssetPatch, NewFolderPost } from "@/lib/models";
import {
  Button,
  Fieldset,
  Grid,
  Group,
  Textarea,
  TextInput,
} from "@mantine/core";
import { useForm } from "@mantine/form";
import { useMutation, useQueryClient } from "@tanstack/react-query";

export function AssetNewForm({
  parent,
  onClose,
}: {
  parent: Asset;
  onClose: () => void;
}) {
  return (
    <>
      <Grid grow>
        <Grid.Col span={{ base: 12, md: 6, lg: 3 }}>
          <Fieldset legend="Upload"></Fieldset>
        </Grid.Col>
        <Grid.Col span={{ base: 12, md: 6, lg: 3 }}>
          <Fieldset legend="New Folder">
            <NewFolder
              parent={parent}
              callback={() => {
                onClose();
              }}
            />
          </Fieldset>
        </Grid.Col>
        <Grid.Col span={{ base: 12, md: 6, lg: 3 }}>
          <Fieldset legend="Import">
            <ImportList
              parent={parent}
              callback={() => {
                onClose();
              }}
            />
          </Fieldset>
        </Grid.Col>
      </Grid>
    </>
  );
}

interface nestedProps {
  parent: Asset;
  callback: () => void;
}

function NewFolder({ parent, callback }: nestedProps) {
  const queryClient = useQueryClient();
  const mutation = useMutation({
    mutationFn: newFolderPost,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["asset", parent.ID] });
      callback();
    },
  });

  const form = useForm<NewFolderPost>({
    initialValues: {
      ParentID: parent.ID,
      FolderName: "",
    } as NewFolderPost,
    validate: {
      FolderName: (value) =>
        !value || value.length > 0 ? undefined : "Folder name is required",
    },
  });

  const formSubmit = (values: AssetPatch): Promise<any> => {
    return mutation.mutateAsync(values);
  };

  return (
    <form onSubmit={form.onSubmit(formSubmit)} onReset={form.reset}>
      <TextInput
        placeholder="Folder name"
        key={form.key("FolderName")}
        {...form.getInputProps("FolderName")}
      />
      <Group align="flex-end" mt="lg">
        <Button type="submit">Create</Button>
      </Group>
    </form>
  );
}

function ImportList({ parent, callback }: nestedProps) {
  return (
    <form>
      <Textarea
        placeholder="List of links"
        minRows={5}
        maxRows={10}
        onChange={(e) => console.log(e.target.value)}
      />
      <Group align="flex-end" mt="lg">
        <Button onClick={callback}>Import</Button>
      </Group>
    </form>
  );
}
