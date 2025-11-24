import React from 'react';
import {
  Button,
  Fieldset,
  Grid,
  Group,
  Textarea,
  TextInput,
} from '@mantine/core';
import { useForm } from '@mantine/form';
import { useMutation, useQueryClient } from '@tanstack/react-query';

import { newFolderPost } from 'lib/fetchers/getAsset';
import { type Asset, type AssetPatch, type NewFolderPost } from 'lib/models';

interface INestedProps {
  readonly parent: Asset;
  readonly callback: () => void;
}

/**
 *
 * @param param0
 * @param param0.parent
 * @param param0.callback
 * @returns
 */
function NewFolder({ parent, callback }: INestedProps) {
  const queryClient = useQueryClient();
  const mutation = useMutation({
    mutationFn: newFolderPost,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['asset', parent.ID] }).catch(console.error);
      callback();
    },
  });

  const form = useForm<NewFolderPost>({
    initialValues: {
      ParentID: parent.ID,
      FolderName: '',
    } as NewFolderPost,
    validate: {
      FolderName: (value) => ((!value || value.length > 0) ? undefined : 'Folder name is required'),
    },
  });

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const formSubmit = (values: AssetPatch): Promise<any> => (
    mutation.mutateAsync(values)
  );

  return (
    <form onSubmit={form.onSubmit(formSubmit)} onReset={form.reset}>
      <TextInput
        placeholder="Folder name"
        key={form.key('FolderName')}
        // eslint-disable-next-line react/jsx-props-no-spreading
        {...form.getInputProps('FolderName')}
      />
      <Group align="flex-end" mt="lg">
        <Button type="submit">Create</Button>
      </Group>
    </form>
  );
}

/**
 *
 * @param param0
 * @param param0.parent
 * @param param0.callback
 * @returns
 */
// eslint-disable-next-line @typescript-eslint/no-unused-vars
function ImportList({ parent, callback }: INestedProps) {
  return (
    <form>
      <Textarea
        placeholder="List of links"
        minRows={5}
        maxRows={10}
        // eslint-disable-next-line no-console
        onChange={(e) => console.log(e.target.value)}
      />
      <Group align="flex-end" mt="lg">
        <Button onClick={callback}>Import</Button>
      </Group>
    </form>
  );
}

/**
 *
 * @param param0
 * @param param0.parent
 * @param param0.onClose
 * @returns
 */
function AssetNewForm({
  parent,
  onClose,
}: {
  readonly parent: Asset;
  readonly onClose: () => void;
}) {
  return (
    <Grid grow>
      <Grid.Col span={{ base: 12, md: 6, lg: 3 }}>
        <Fieldset legend="Upload" />
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
  );
}

export default AssetNewForm;
