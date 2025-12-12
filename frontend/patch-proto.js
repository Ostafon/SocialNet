const fs = require('fs');
const path = require('path');

console.log('🚨 EMERGENCY PATCH - Fixing proto files...\n');

const protoDir = path.join(__dirname, 'src/chat/proto');

function patchFile(filePath) {
  if (!fs.existsSync(filePath)) {
    console.log(`⚠️  Not found: ${filePath}`);
    return false;
  }

  let content = fs.readFileSync(filePath, 'utf8');
  const original = content;

  // Заменяем ВСЕ проблемные методы
  content = content.replace(
    /reader\.readStringRequireUtf8\(\)/g,
    'reader.readString()'
  );
  content = content.replace(/\.readStringRequireUtf8\(\)/g, '.readString()');

  // Если файл изменился - сохраняем
  if (content !== original) {
    fs.writeFileSync(filePath, content, 'utf8');
    console.log(`✅ PATCHED: ${path.basename(filePath)}`);

    // Показываем количество замен
    const count = (original.match(/readStringRequireUtf8/g) || []).length;
    console.log(`   Replaced ${count} occurrences\n`);
    return true;
  } else {
    console.log(`ℹ️  No changes: ${path.basename(filePath)}\n`);
    return false;
  }
}

// Патчим все JS файлы
if (fs.existsSync(protoDir)) {
  const files = fs.readdirSync(protoDir).filter(f => f.endsWith('.js'));

  console.log(`📂 Found ${files.length} files in ${protoDir}\n`);

  let patchedCount = 0;
  files.forEach(file => {
    if (patchFile(path.join(protoDir, file))) {
      patchedCount++;
    }
  });

  console.log('═══════════════════════════════════════');
  console.log(`🎉 SUCCESS! Patched ${patchedCount}/${files.length} files`);
  console.log('═══════════════════════════════════════\n');
  console.log('🔄 Now restart your dev server:');
  console.log('   npm start');
} else {
  console.log(`❌ ERROR: Directory not found!`);
  console.log(`   Expected: ${protoDir}`);
  console.log(`\n💡 Make sure your proto files are in: src/chat/proto/`);
}
