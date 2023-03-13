add_rules("mode.debug", "mode.release")

set_arch("x64")
set_kind("binary")

set_languages("c++11")

for _, fi in ipairs(os.files("*.c")) do
  target(path.basename(fi))
    add_files(fi)
end
