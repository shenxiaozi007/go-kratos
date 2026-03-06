const fs = require('fs');
const path = require('path');
const {
  Document,
  Packer,
  Paragraph,
  TextRun,
  HeadingLevel,
  AlignmentType,
  convertInchesToTwip,
} = require('docx');

const readmePath = path.join(__dirname, '..', 'README.md');
const outPath = path.join(__dirname, '..', 'README.docx');

const raw = fs.readFileSync(readmePath, 'utf8');

// Simple markdown parse: split by lines, detect ## and ``` blocks
const lines = raw.split(/\r?\n/);
const children = [];
let inCode = false;
let codeLines = [];
const flushCode = () => {
  if (codeLines.length === 0) return;
  const code = codeLines.join('\n');
  codeLines = [];
  children.push(
    new Paragraph({
      children: [
        new TextRun({
          text: code,
          font: 'Consolas',
          size: 20,
        }),
      ],
    })
  );
};

for (let i = 0; i < lines.length; i++) {
  const line = lines[i];
  if (line.startsWith('```')) {
    if (inCode) {
      flushCode();
      inCode = false;
    } else {
      inCode = true;
    }
    continue;
  }
  if (inCode) {
    codeLines.push(line);
    continue;
  }
  if (line.startsWith('# ')) {
    children.push(
      new Paragraph({
        heading: HeadingLevel.HEADING_1,
        children: [new TextRun({ text: line.slice(2).trim(), bold: true })],
      })
    );
    continue;
  }
  if (line.startsWith('## ')) {
    children.push(
      new Paragraph({
        heading: HeadingLevel.HEADING_2,
        children: [new TextRun({ text: line.slice(3).trim(), bold: true })],
      })
    );
    continue;
  }
  if (line.trim() === '') {
    children.push(new Paragraph({ text: '' }));
    continue;
  }
  children.push(
    new Paragraph({
      children: [new TextRun({ text: line })],
    })
  );
}
flushCode();

const doc = new Document({
  sections: [
    {
      properties: {
        page: {
          margin: {
            top: convertInchesToTwip(1),
            right: convertInchesToTwip(1),
            bottom: convertInchesToTwip(1),
            left: convertInchesToTwip(1),
          },
        },
      },
      children,
    },
  ],
});

Packer.toBuffer(doc).then((buffer) => {
  fs.writeFileSync(outPath, buffer);
  console.log('Written:', outPath);
});
