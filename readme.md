# msstg

[![Go checker](https://github.com/alekssamos/msstg/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/alekssamos/msstg/actions/workflows/go.yml)

## Description
msstg is a Telegram bot for converting text to speech. It is completely free with no limits on the size and number of characters.

## Roadmap
- [ ] Multilingual user interface for the bot
- [x] Voice synthesis for short messages
- [ ] Conversion from mp3 to ogg
- [x] Voice selection via buttons
- [ ] Voice selection via commands
- [ ] Voice pitch selection via buttons
- [ ] Voice pitch selection via commands
- [ ] Voice rate selection via buttons
- [ ] Voice rate selection via commands
- [ ] Support for text files in various formats (txt, docx, fb2, epub)
- [ ] Splitting books into chapters with configurable parameters
- [ ] Sending audio chapters of books as a zip archive instead of individual files with configurable parameters

## Author
[@alekssamos](https://github.com/alekssamos)

## Running and Building
Before running, export the environment variable `BOT_TOKEN` obtained from Telegram [Bot Father](https://t.me/BotFather).

You can download the ready builds from the GitHub releases page.

### Build
```bash
go mod tidy
go build .
```

## Contribute
We welcome contributions! Here’s how you can get started:

1. **Fork the Repository**
   Click the "Fork" button at the top right of this page to create a copy of the repo under your GitHub account.

2. **Clone the Repository**
   Use the following command to clone your forked repository to your local machine:
   ```bash
   git clone https://github.com/YOUR_USERNAME/msstg.git
   ```

3. Install golang 1.24.1+
   Download Go version 1.24.1 or a later   from [official website](https://go.dev/dl/) or use system package manager.

4. Install pre-commit
   Make sure you have Python and pip installed, then run:

   ```bash
   pip install pre-commit
   ```

   After installing, set it up in your repository:

   ```bash
   pre-commit install
   ```

5. Create a New Branch
   Create a new branch for your changes:

   ```bash
   git checkout -b your-feature-branch
   ```

6. Make Changes
   Edit the code in your favorite text editor and make the desired changes.

7. Commit Your Changes
   Stage your changes and commit them with a descriptive message:

   ```bash
   git add .
   git commit -m "Description of your changes"
   ```

8. Push to the Remote Repository
   Push your changes to your forked repository:

   ```bash
   git push origin your-feature-branch
   ```

9. Create a Pull Request
   Go to your fork on GitHub, click on "Pull requests," and then "New pull request." Follow the prompts to submit your changes for review.

## Notes
This is a rewritten MS Speech Bot from Python to golang, which everyone already knows, written in 2021.
I did not publish the python code of the bot here, as it is terrible.

---

© 2025 @alekssamos
