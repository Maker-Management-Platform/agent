import React from 'react';
import {
  Button,
  Fieldset,
  Grid,
  Group,
  Textarea,
  TextInput,
} from '@mantine/core';

import { useNewFolderPost } from 'fetchers/getAsset';
import { type Asset } from 'types/Asset';

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
  const [folderName, setFolderName] = React.useState('');
  const postNewFolder = useNewFolderPost();

  const submitFolderName = () => {
    postNewFolder({
      ParentID: parent.ID,
      FolderName: folderName,
    }).then(() => callback()).catch(console.error);
  };

  return (
    <div>
      <TextInput
        onChange={(e) => setFolderName(e.target.value.trim())}
        defaultValue={folderName}
        placeholder="Folder name"
      />
      <Group align="flex-end" mt="lg">
        <Button onClick={submitFolderName} disabled={folderName.length === 0}>Create</Button>
      </Group>
    </div>
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
        // onChange={(e) => console.log(parent.ID, e.target.value)}
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
