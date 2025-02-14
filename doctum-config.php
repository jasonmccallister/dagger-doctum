<?php

use Doctum\Doctum;
use Symfony\Component\Finder\Finder;

$iterator = Finder::create()
    ->files()
    ->name("*.php")
    ->exclude(".changes")
    ->exclude("docker")
    ->exclude("runtime")
    ->exclude("tests")
    ->exclude("src/Codegen/")
    ->exclude("src/Command/")
    ->exclude("src/Connection/")
    ->exclude("src/Exception/")
    ->exclude("src/GraphQl/")
    ->exclude("src/Service/")
    ->exclude("src/ValueObject/")
    ->in("/work/repository/");

return new Doctum($iterator);
