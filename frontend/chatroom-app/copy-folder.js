import fs from 'fs-extra';

// Define source and destination paths
const sourceFolder = 'src/assets';
const destinationFolder = 'dist/src/assets';

// Copy the folder
fs.copy(sourceFolder, destinationFolder, err => {
    if (err) {
        console.error('Error copying folder:', err);
    } else {
        console.log('Folder copied successfully!');
    }
});