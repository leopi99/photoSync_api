-- phpMyAdmin SQL Dump
-- version 5.0.2
-- https://www.phpmyadmin.net/
--
-- Host: 127.0.0.1
-- Creato il: Giu 19, 2021 alle 21:06
-- Versione del server: 10.4.11-MariaDB
-- Versione PHP: 7.4.4

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Database: `photosync`
--

-- --------------------------------------------------------

--
-- Struttura della tabella `object`
--

CREATE TABLE `object` (
  `objectID` int(11) NOT NULL,
  `userID` int(11) NOT NULL,
  `local_path` text NOT NULL,
  `creation_date` varchar(30) NOT NULL,
  `sync_date` varchar(30) NOT NULL,
  `picture_position` text NOT NULL,
  `type` varchar(12) NOT NULL,
  `byte_size` bigint(20) UNSIGNED NOT NULL,
  `extension` varchar(6) NOT NULL,
  `local_id` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- --------------------------------------------------------

--
-- Struttura della tabella `user`
--

CREATE TABLE `user` (
  `userID` int(11) NOT NULL,
  `username` varchar(15) NOT NULL,
  `password` text NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

--
-- Indici per le tabelle scaricate
--

--
-- Indici per le tabelle `object`
--
ALTER TABLE `object`
  ADD PRIMARY KEY (`objectID`),
  ADD KEY `fk_foreign_userID` (`userID`);

--
-- Indici per le tabelle `user`
--
ALTER TABLE `user`
  ADD PRIMARY KEY (`userID`);

--
-- AUTO_INCREMENT per le tabelle scaricate
--

--
-- AUTO_INCREMENT per la tabella `object`
--
ALTER TABLE `object`
  MODIFY `objectID` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT per la tabella `user`
--
ALTER TABLE `user`
  MODIFY `userID` int(11) NOT NULL AUTO_INCREMENT;

--
-- Limiti per le tabelle scaricate
--

--
-- Limiti per la tabella `object`
--
ALTER TABLE `object`
  ADD CONSTRAINT `fk_foreign_userID` FOREIGN KEY (`userID`) REFERENCES `user` (`userID`) ON DELETE CASCADE ON UPDATE CASCADE;
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
