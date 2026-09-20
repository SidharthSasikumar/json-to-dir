Running the Command
Once everything is set up, you can build and run the CLI tool:

Build the project:

```sh
go build -o json-cli
```
Run the process command:

```sh

./json-cli process input.json output_directory
```
Replace input.json with the path to your input JSON file and output_directory with the directory where you want the output to be saved.
Requires Go 1.25 or later. Output uses rooted filesystem operations to prevent JSON keys or symlinks from writing outside the output directory.
